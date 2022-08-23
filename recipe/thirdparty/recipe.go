/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package thirdparty

import (
	"errors"
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "thirdparty"

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       tpmodels.TypeNormalisedInput
	RecipeImpl   tpmodels.RecipeInterface
	APIImpl      tpmodels.APIInterface
	Providers    []tpmodels.TypeProvider
}

var singletonInstance *Recipe

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *tpmodels.TypeInput, emailDeliveryIngredient *emaildelivery.Ingredient, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (Recipe, error) {
	r := &Recipe{}

	r.RecipeModule = supertokens.MakeRecipeModule(recipeId, appInfo, r.handleAPIRequest, r.getAllCORSHeaders, r.getAPIsHandled, r.handleError, onSuperTokensAPIError)

	querierInstance, err := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	if err != nil {
		return Recipe{}, err
	}
	verifiedConfig, err := validateAndNormaliseUserInput(r, appInfo, config)
	if err != nil {
		return Recipe{}, err
	}
	r.Config = verifiedConfig
	r.APIImpl = verifiedConfig.Override.APIs(api.MakeAPIImplementation())
	r.RecipeImpl = verifiedConfig.Override.Functions(MakeRecipeImplementation(*querierInstance))
	r.Providers = config.SignInAndUpFeature.Providers

	return *r, nil
}

func recipeInit(config *tpmodels.TypeInput) supertokens.Recipe {
	return func(appInfo supertokens.NormalisedAppinfo, onSuperTokensAPIError func(err error, req *http.Request, res http.ResponseWriter)) (*supertokens.RecipeModule, error) {
		if singletonInstance == nil {
			recipe, err := MakeRecipe(RECIPE_ID, appInfo, config, nil, onSuperTokensAPIError)
			if err != nil {
				return nil, err
			}
			singletonInstance = &recipe
			return &singletonInstance.RecipeModule, nil
		}
		return nil, errors.New("ThirdParty recipe has already been initialised. Please check your code for bugs.")
	}
}

func getRecipeInstanceOrThrowError() (*Recipe, error) {
	if singletonInstance != nil {
		return singletonInstance, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the init function?")
}

// implement RecipeModule

func (r *Recipe) getAPIsHandled() ([]supertokens.APIHandled, error) {
	signInUpAPI, err := supertokens.NewNormalisedURLPath(SignInUpAPI)
	if err != nil {
		return nil, err
	}
	authorisationAPI, err := supertokens.NewNormalisedURLPath(AuthorisationAPI)
	if err != nil {
		return nil, err
	}
	appleRedirectHandlerAPI, err := supertokens.NewNormalisedURLPath(AppleRedirectHandlerAPI)
	if err != nil {
		return nil, err
	}
	return append([]supertokens.APIHandled{{
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: signInUpAPI,
		ID:                     SignInUpAPI,
		Disabled:               r.APIImpl.SignInUpPOST == nil,
	}, {
		Method:                 http.MethodGet,
		PathWithoutAPIBasePath: authorisationAPI,
		ID:                     AuthorisationAPI,
		Disabled:               r.APIImpl.AuthorisationUrlGET == nil,
	}, {
		Method:                 http.MethodPost,
		PathWithoutAPIBasePath: appleRedirectHandlerAPI,
		ID:                     AppleRedirectHandlerAPI,
		Disabled:               r.APIImpl.AppleRedirectHandlerPOST == nil,
	}}), nil
}

func (r *Recipe) handleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, path supertokens.NormalisedURLPath, method string) error {
	options := tpmodels.APIOptions{
		Config:               r.Config,
		OtherHandler:         theirHandler,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Providers:            r.Providers,
		Req:                  req,
		Res:                  res,
		AppInfo:              r.RecipeModule.GetAppInfo(),
	}
	if id == SignInUpAPI {
		return api.SignInUpAPI(r.APIImpl, options)
	} else if id == AuthorisationAPI {
		return api.AuthorisationUrlAPI(r.APIImpl, options)
	} else if id == AppleRedirectHandlerAPI {
		return api.AppleRedirectHandler(r.APIImpl, options)
	}
	return errors.New("should never come here")
}

func (r *Recipe) getAllCORSHeaders() []string {
	return []string{}
}

func (r *Recipe) handleError(err error, req *http.Request, res http.ResponseWriter) (bool, error) {
	return false, err
}

func (r *Recipe) getEmailForUserId(userID string, userContext supertokens.UserContext) (string, error) {
	userInfo, err := (*r.RecipeImpl.GetUserByID)(userID, userContext)
	if err != nil {
		return "", err
	}
	if userInfo == nil {
		return "", errors.New("unknown User ID provided")
	}
	return userInfo.Email, nil
}

func ResetForTest() {
	singletonInstance = nil
}
