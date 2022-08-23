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

package api

import (
	"errors"
	"fmt"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/claims"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() evmodels.APIInterface {
	verifyEmailPOST := func(token string, options evmodels.APIOptions, sessionContainer *sessmodels.SessionContainer, userContext supertokens.UserContext) (evmodels.VerifyEmailPOSTResponse, error) {
		resp, err := (*options.RecipeImplementation.VerifyEmailUsingToken)(token, userContext)
		if err != nil {
			return evmodels.VerifyEmailPOSTResponse{}, err
		}
		if resp.OK != nil {
			if sessionContainer != nil {
				sessionContainer.FetchAndSetClaim(claims.EmailVerificationClaim.TypeSessionClaim, userContext)
			}
			return evmodels.VerifyEmailPOSTResponse{
				OK: resp.OK,
			}, err
		} else {
			return evmodels.VerifyEmailPOSTResponse{
				EmailVerificationInvalidTokenError: resp.EmailVerificationInvalidTokenError,
			}, nil
		}
	}

	isEmailVerifiedGET := func(options evmodels.APIOptions, sessionContainer *sessmodels.SessionContainer, userContext supertokens.UserContext) (evmodels.IsEmailVerifiedGETResponse, error) {
		if sessionContainer == nil {
			return evmodels.IsEmailVerifiedGETResponse{}, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
		}

		sessionContainer.FetchAndSetClaim(claims.EmailVerificationClaim.TypeSessionClaim, userContext)
		isVerified, err := sessionContainer.GetClaimValue(claims.EmailVerificationClaim.TypeSessionClaim, userContext)
		if err != nil {
			return evmodels.IsEmailVerifiedGETResponse{}, err
		}

		if err != nil {
			return evmodels.IsEmailVerifiedGETResponse{}, err
		}
		if err != nil {
			return evmodels.IsEmailVerifiedGETResponse{}, err
		}
		return evmodels.IsEmailVerifiedGETResponse{
			OK: &struct{ IsVerified bool }{
				IsVerified: isVerified.(bool),
			},
		}, nil
	}

	generateEmailVerifyTokenPOST := func(options evmodels.APIOptions, sessionContainer *sessmodels.SessionContainer, userContext supertokens.UserContext) (evmodels.GenerateEmailVerifyTokenPOSTResponse, error) {
		stInstance, err := supertokens.GetInstanceOrThrowError()
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}
		if sessionContainer == nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, supertokens.BadInputError{Msg: "Session is undefined. Should not come here."}
		}

		userID := sessionContainer.GetUserIDWithContext(userContext)
		email, err := options.GetEmailForUserID(userID, userContext)
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}
		if email.UnknownUserIDError != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, errors.New("unknown userid")
		}
		if email.EmailDoesNotExistError != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{
				EmailAlreadyVerifiedError: &struct{}{},
			}, nil
		}
		response, err := (*options.RecipeImplementation.CreateEmailVerificationToken)(userID, email.OK.Email, userContext)
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}

		if response.EmailAlreadyVerifiedError != nil {
			supertokens.LogDebugMessage(fmt.Sprintf("Email verification email not sent to %s because it is already verified", email.OK.Email))
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{
				EmailAlreadyVerifiedError: &struct{}{},
			}, nil
		}

		user := evmodels.User{
			ID:    userID,
			Email: email.OK.Email,
		}
		emailVerificationURL := fmt.Sprintf(
			"%s%s/verify-email?token=%s&rid=%s",
			stInstance.AppInfo.WebsiteDomain.GetAsStringDangerous(),
			stInstance.AppInfo.WebsiteBasePath.GetAsStringDangerous(),
			response.OK.Token,
			options.RecipeID,
		)

		supertokens.LogDebugMessage(fmt.Sprintf("Sending email verification email to %s", email.OK.Email))
		err = (*options.EmailDelivery.IngredientInterfaceImpl.SendEmail)(emaildelivery.EmailType{
			EmailVerification: &emaildelivery.EmailVerificationType{
				User: emaildelivery.User{
					ID:    user.ID,
					Email: user.Email,
				},
				EmailVerifyLink: emailVerificationURL,
			},
		}, userContext)
		if err != nil {
			return evmodels.GenerateEmailVerifyTokenPOSTResponse{}, err
		}

		return evmodels.GenerateEmailVerifyTokenPOSTResponse{
			OK: &struct{}{},
		}, nil
	}

	return evmodels.APIInterface{
		VerifyEmailPOST:              &verifyEmailPOST,
		IsEmailVerifiedGET:           &isEmailVerifiedGET,
		GenerateEmailVerifyTokenPOST: &generateEmailVerifyTokenPOST,
	}
}
