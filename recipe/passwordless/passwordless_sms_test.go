/*
 * Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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

package passwordless

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/smsdelivery/twilioService"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestSmsBackwardCompatibilityServiceForContactPhoneMethod(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				GetCustomUserInputCode: func(userContext supertokens.UserContext) (string, error) {
					return "", nil
				},
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	someCode := "someCode"
	(*singletonInstance.SmsDelivery.IngredientInterfaceImpl.SendSms)(smsdelivery.SmsType{
		PasswordlessLogin: &smsdelivery.PasswordlessLoginType{
			PhoneNumber:      "somePhoneNumber",
			UserInputCode:    &someCode,
			UrlWithLinkCode:  nil,
			CodeLifetime:     3600,
			PreAuthSessionId: "someSession",
		},
	}, nil)

	assert.Equal(t, PasswordlessLoginSmsSentForTest, true)
}

func TestSmsBackwardCompatibilityServiceForContactEmailOrPhoneMethod(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				GetCustomUserInputCode: func(userContext supertokens.UserContext) (string, error) {
					return "", nil
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	someCode := "someCode"
	(*singletonInstance.SmsDelivery.IngredientInterfaceImpl.SendSms)(smsdelivery.SmsType{
		PasswordlessLogin: &smsdelivery.PasswordlessLoginType{
			PhoneNumber:      "somePhoneNumber",
			PreAuthSessionId: "someSession",
			UserInputCode:    &someCode,
			UrlWithLinkCode:  nil,
			CodeLifetime:     3600,
		},
	}, nil)

	assert.Equal(t, PasswordlessLoginSmsSentForTest, true)
}

func TestSmsBackwardCompatibilityServiceWithtCustomFunctionForContactPhoneMethod(t *testing.T) {
	customSmsSent := false

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				GetCustomUserInputCode: func(userContext supertokens.UserContext) (string, error) {
					return "", nil
				},
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						customSmsSent = true
						return nil
					},
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	someCode := "someCode"
	(*singletonInstance.SmsDelivery.IngredientInterfaceImpl.SendSms)(smsdelivery.SmsType{
		PasswordlessLogin: &smsdelivery.PasswordlessLoginType{
			PhoneNumber:      "somePhoneNumber",
			UserInputCode:    &someCode,
			UrlWithLinkCode:  nil,
			CodeLifetime:     3600,
			PreAuthSessionId: "someSession",
		},
	}, nil)

	assert.Equal(t, PasswordlessLoginSmsSentForTest, false)
	assert.Equal(t, customSmsSent, true)
}

func TestSmsBackwardCompatibilityServiceWithtCustomFunctionForContactEmailOrPhoneMethod(t *testing.T) {
	customSmsSent := false

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				GetCustomUserInputCode: func(userContext supertokens.UserContext) (string, error) {
					return "", nil
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						customSmsSent = true
						return nil
					},
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	someCode := "someCode"
	(*singletonInstance.SmsDelivery.IngredientInterfaceImpl.SendSms)(smsdelivery.SmsType{
		PasswordlessLogin: &smsdelivery.PasswordlessLoginType{
			PhoneNumber:      "somePhoneNumber",
			PreAuthSessionId: "someSession",
			UserInputCode:    &someCode,
			UrlWithLinkCode:  nil,
			CodeLifetime:     3600,
		},
	}, nil)

	assert.Equal(t, PasswordlessLoginSmsSentForTest, false)
	assert.Equal(t, customSmsSent, true)
}

func TestSmsBackwardCompatibilityServiceWithOverrideForContactPhoneMethod(t *testing.T) {
	funcCalled := false
	overrideCalled := false

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						funcCalled = true
						return nil
					},
				},
				SmsDelivery: &smsdelivery.TypeInput{
					Override: func(originalImplementation smsdelivery.SmsDeliveryInterface) smsdelivery.SmsDeliveryInterface {
						(*originalImplementation.SendSms) = func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
							overrideCalled = true
							return nil
						}
						return originalImplementation
					},
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	someCode := "someCode"
	(*singletonInstance.SmsDelivery.IngredientInterfaceImpl.SendSms)(smsdelivery.SmsType{
		PasswordlessLogin: &smsdelivery.PasswordlessLoginType{
			PhoneNumber:      "somePhoneNumber",
			PreAuthSessionId: "someSession",
			UserInputCode:    &someCode,
			UrlWithLinkCode:  nil,
			CodeLifetime:     3600,
		},
	}, nil)

	assert.Equal(t, PasswordlessLoginSmsSentForTest, false)
	assert.Equal(t, funcCalled, false)
	assert.Equal(t, overrideCalled, true)
}

func TestSmsBackwardCompatibilityServiceWithOverrideForContactEmailOrPhoneMethod(t *testing.T) {
	funcCalled := false
	overrideCalled := false

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						funcCalled = true
						return nil
					},
				},
				SmsDelivery: &smsdelivery.TypeInput{
					Override: func(originalImplementation smsdelivery.SmsDeliveryInterface) smsdelivery.SmsDeliveryInterface {
						(*originalImplementation.SendSms) = func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
							overrideCalled = true
							return nil
						}
						return originalImplementation
					},
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	someCode := "someCode"
	(*singletonInstance.SmsDelivery.IngredientInterfaceImpl.SendSms)(smsdelivery.SmsType{
		PasswordlessLogin: &smsdelivery.PasswordlessLoginType{
			PhoneNumber:      "somePhoneNumber",
			PreAuthSessionId: "someSession",
			UserInputCode:    &someCode,
			UrlWithLinkCode:  nil,
			CodeLifetime:     3600,
		},
	}, nil)

	assert.Equal(t, PasswordlessLoginSmsSentForTest, false)
	assert.Equal(t, funcCalled, false)
	assert.Equal(t, overrideCalled, true)
}

func TestTwilioServiceOverrideForContactPhoneMethod(t *testing.T) {
	getContentCalled := false
	sendRawSmsCalled := false
	customCalled := false

	fromPhoneNumber := "someNumber"
	twilioService := twilioService.MakeTwilioService(
		smsdelivery.TwilioTypeInput{
			TwilioSettings: smsdelivery.TwilioServiceConfig{
				AccountSid:          "sid",
				AuthToken:           "token",
				From:                &fromPhoneNumber,
				MessagingServiceSid: nil,
			},
			Override: func(originalImplementation smsdelivery.TwilioServiceInterface) smsdelivery.TwilioServiceInterface {
				(*originalImplementation.GetContent) = func(input smsdelivery.SmsType, userContext supertokens.UserContext) (smsdelivery.GetContentResult, error) {
					getContentCalled = true
					return smsdelivery.GetContentResult{}, nil
				}

				(*originalImplementation.SendRawSms) = func(input smsdelivery.GetContentResult, userContext supertokens.UserContext) error {
					sendRawSmsCalled = true
					return nil
				}

				return originalImplementation
			},
		},
	)

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &twilioService,
				},
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						customCalled = true
						return nil
					},
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	someCode := "someCode"
	(*singletonInstance.SmsDelivery.IngredientInterfaceImpl.SendSms)(smsdelivery.SmsType{
		PasswordlessLogin: &smsdelivery.PasswordlessLoginType{
			PhoneNumber:      "somePhoneNumber",
			PreAuthSessionId: "someSession",
			UserInputCode:    &someCode,
			UrlWithLinkCode:  nil,
			CodeLifetime:     3600,
		},
	}, nil)

	assert.Nil(t, err)
	assert.Equal(t, customCalled, false)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawSmsCalled, true)
}

func TestTwilioServiceOverrideForContactEmailOrPhoneMethod(t *testing.T) {
	getContentCalled := false
	sendRawSmsCalled := false
	customCalled := false

	fromPhoneNumber := "someNumber"
	twilioService := twilioService.MakeTwilioService(
		smsdelivery.TwilioTypeInput{
			TwilioSettings: smsdelivery.TwilioServiceConfig{
				AccountSid:          "sid",
				AuthToken:           "token",
				From:                &fromPhoneNumber,
				MessagingServiceSid: nil,
			},
			Override: func(originalImplementation smsdelivery.TwilioServiceInterface) smsdelivery.TwilioServiceInterface {
				(*originalImplementation.GetContent) = func(input smsdelivery.SmsType, userContext supertokens.UserContext) (smsdelivery.GetContentResult, error) {
					getContentCalled = true
					return smsdelivery.GetContentResult{}, nil
				}

				(*originalImplementation.SendRawSms) = func(input smsdelivery.GetContentResult, userContext supertokens.UserContext) error {
					sendRawSmsCalled = true
					return nil
				}

				return originalImplementation
			},
		},
	)

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",
				SmsDelivery: &smsdelivery.TypeInput{
					Service: &twilioService,
				},
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						customCalled = true
						return nil
					},
				},
			}),
		},
	}

	BeforeEach()
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	someCode := "someCode"
	(*singletonInstance.SmsDelivery.IngredientInterfaceImpl.SendSms)(smsdelivery.SmsType{
		PasswordlessLogin: &smsdelivery.PasswordlessLoginType{
			PhoneNumber:      "somePhoneNumber",
			PreAuthSessionId: "someSession",
			UserInputCode:    &someCode,
			UrlWithLinkCode:  nil,
			CodeLifetime:     3600,
		},
	}, nil)

	assert.Nil(t, err)
	assert.Equal(t, customCalled, false)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawSmsCalled, true)
}

func TestTwilioServiceOverrideForContactPhoneMethodThroughAPI(t *testing.T) {
	getContentCalled := false
	sendRawSmsCalled := false
	customCalled := false

	fromPhoneNumber := "someNumber"
	twilioService := twilioService.MakeTwilioService(
		smsdelivery.TwilioTypeInput{
			TwilioSettings: smsdelivery.TwilioServiceConfig{
				AccountSid:          "sid",
				AuthToken:           "token",
				From:                &fromPhoneNumber,
				MessagingServiceSid: nil,
			},
			Override: func(originalImplementation smsdelivery.TwilioServiceInterface) smsdelivery.TwilioServiceInterface {
				(*originalImplementation.GetContent) = func(input smsdelivery.SmsType, userContext supertokens.UserContext) (smsdelivery.GetContentResult, error) {
					getContentCalled = true
					return smsdelivery.GetContentResult{}, nil
				}

				(*originalImplementation.SendRawSms) = func(input smsdelivery.GetContentResult, userContext supertokens.UserContext) error {
					sendRawSmsCalled = true
					return nil
				}

				return originalImplementation
			},
		},
	)

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE",

				SmsDelivery: &smsdelivery.TypeInput{
					Service: &twilioService,
				},
				ContactMethodPhone: plessmodels.ContactMethodPhoneConfig{
					Enabled: true,
					CreateAndSendCustomTextMessage: func(phoneNumber string, userInputCode, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
						customCalled = true
						return nil
					},
					ValidatePhoneNumber: func(phoneNumber interface{}) *string {
						return nil
					},
				},
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	unittesting.PasswordlessPhoneLoginRequest("somePhone", testServer.URL)

	assert.Equal(t, customCalled, false)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawSmsCalled, true)
}
