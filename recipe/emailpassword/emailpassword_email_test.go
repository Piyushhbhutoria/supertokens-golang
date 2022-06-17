/*
 * Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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

package emailpassword

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultBackwardCompatibilityPasswordResetForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	testServer := supertokensInitForTest(t, session.Init(nil), Init(nil))
	defer testServer.Close()

	SignUp("test@example.com", "1234abcd")
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, PasswordResetEmailSentForTest)
	assert.Equal(t, PasswordResetDataForTest.User.Email, "test@example.com")
	assert.NotEmpty(t, PasswordResetDataForTest.PasswordResetURLWithToken)
}

func TestDefaultBackwardCompatibilityPasswordResetForNonExistantUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	testServer := supertokensInitForTest(t, session.Init(nil), Init(nil))
	defer testServer.Close()

	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, PasswordResetEmailSentForTest)
	assert.Empty(t, PasswordResetDataForTest.User.Email)
	assert.Empty(t, PasswordResetDataForTest.PasswordResetURLWithToken)
}

func TestBackwardCompatibilityResetPasswordForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	passwordResetLink := ""

	tpepConfig := &epmodels.TypeInput{
		ResetPasswordUsingTokenFeature: &epmodels.TypeInputResetPasswordUsingTokenFeature{
			CreateAndSendCustomEmail: func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext) {
				email = user.Email
				passwordResetLink = passwordResetURLWithToken
				customCalled = true
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	SignUp("test@example.com", "1234abcd")
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordResetEmailSentForTest)
	assert.Empty(t, PasswordResetDataForTest.User.Email)
	assert.Empty(t, PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler called
	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, passwordResetLink)
	assert.True(t, customCalled)
}

func TestBackwardCompatibilityResetPasswordForNonExistantUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	passwordResetLink := ""

	tpepConfig := &epmodels.TypeInput{
		ResetPasswordUsingTokenFeature: &epmodels.TypeInputResetPasswordUsingTokenFeature{
			CreateAndSendCustomEmail: func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext) {
				email = user.Email
				passwordResetLink = passwordResetURLWithToken
				customCalled = true
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordResetEmailSentForTest)
	assert.Empty(t, PasswordResetDataForTest.User.Email)
	assert.Empty(t, PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, passwordResetLink)
	assert.False(t, customCalled)
}

func TestCustomOverrideResetPasswordForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	passwordResetLink := ""

	tpepConfig := &epmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
				*originalImplementation.SendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
					if input.PasswordReset != nil {
						customCalled = true
						email = input.PasswordReset.User.Email
						passwordResetLink = input.PasswordReset.PasswordResetLink
					}
					return nil
				}
				return originalImplementation
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	SignUp("test@example.com", "1234abcd")
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordResetEmailSentForTest)
	assert.Empty(t, PasswordResetDataForTest.User.Email)
	assert.Empty(t, PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler called
	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, passwordResetLink)
	assert.True(t, customCalled)
}

func TestCustomOverrideResetPasswordForNonExistantUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	passwordResetLink := ""

	tpepConfig := &epmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
				*originalImplementation.SendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
					if input.PasswordReset != nil {
						customCalled = true
						email = input.PasswordReset.User.Email
						passwordResetLink = input.PasswordReset.PasswordResetLink
					}
					return nil
				}
				return originalImplementation
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordResetEmailSentForTest)
	assert.Empty(t, PasswordResetDataForTest.User.Email)
	assert.Empty(t, PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, passwordResetLink)
	assert.False(t, customCalled)
}

func TestSMTPOverridePasswordResetForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	email := ""
	passwordResetLink := ""

	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
				if input.PasswordReset != nil {
					email = input.PasswordReset.User.Email
					passwordResetLink = input.PasswordReset.PasswordResetLink
					getContentCalled = true
				}
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
	tpepConfig := &epmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Service: &smtpService,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	SignUp("test@example.com", "1234abcd")
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordResetEmailSentForTest)
	assert.Empty(t, PasswordResetDataForTest.User.Email)
	assert.Empty(t, PasswordResetDataForTest.PasswordResetURLWithToken)

	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, passwordResetLink)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestSMTPOverridePasswordResetForNonExistantUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	email := ""
	passwordResetLink := ""

	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
				if input.PasswordReset != nil {
					email = input.PasswordReset.User.Email
					passwordResetLink = input.PasswordReset.PasswordResetLink
					getContentCalled = true
				}
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
	tpepConfig := &epmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Service: &smtpService,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordResetEmailSentForTest)
	assert.Empty(t, PasswordResetDataForTest.User.Email)
	assert.Empty(t, PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, passwordResetLink)
	assert.False(t, getContentCalled)
	assert.False(t, sendRawEmailCalled)
}

func TestDefaultBackwardCompatibilityEmailVerifyForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	testServer := supertokensInitForTest(t, session.Init(nil), Init(nil))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("test@example.com", "1234abcd", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, emailverification.EmailVerificationEmailSentForTest)
	assert.Equal(t, emailverification.EmailVerificationDataForTest.User.Email, "test@example.com")
	assert.NotEmpty(t, emailverification.EmailVerificationDataForTest.EmailVerifyURLWithToken)
}

func TestBackwardCompatibilityEmailVerifyForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""

	tpepConfig := &epmodels.TypeInput{
		EmailVerificationFeature: &epmodels.TypeInputEmailVerificationFeature{
			CreateAndSendCustomEmail: func(user epmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
				email = user.Email
				emailVerifyLink = emailVerificationURLWithToken
				customCalled = true
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("test@example.com", "1234abcd", testServer.URL)
	assert.NoError(t, err)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordResetEmailSentForTest)
	assert.Empty(t, PasswordResetDataForTest.User.Email)
	assert.Empty(t, PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler called
	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, emailVerifyLink)
	assert.True(t, customCalled)
}

func TestCustomOverrideEmailVerifyForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""

	tpepConfig := &epmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
				*originalImplementation.SendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
					if input.EmailVerification != nil {
						customCalled = true
						email = input.EmailVerification.User.Email
						emailVerifyLink = input.EmailVerification.EmailVerifyLink
					}
					return nil
				}
				return originalImplementation
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("test@example.com", "1234abcd", testServer.URL)
	assert.NoError(t, err)
	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordResetEmailSentForTest)
	assert.Empty(t, PasswordResetDataForTest.User.Email)
	assert.Empty(t, PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler called
	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, emailVerifyLink)
	assert.True(t, customCalled)
}

func TestSMTPOverrideEmailVerifyForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	email := ""
	emailVerifyLink := ""

	smtpService := smtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: emaildelivery.SMTPServiceConfig{
			Host: "",
			From: emaildelivery.SMTPServiceFromConfig{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
				if input.EmailVerification != nil {
					email = input.EmailVerification.User.Email
					emailVerifyLink = input.EmailVerification.EmailVerifyLink
					getContentCalled = true
				}
				return emaildelivery.SMTPGetContentResult{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
	tpepConfig := &epmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Service: &smtpService,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("test@example.com", "1234abcd", testServer.URL)
	assert.NoError(t, err)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, PasswordResetEmailSentForTest)
	assert.Empty(t, PasswordResetDataForTest.User.Email)
	assert.Empty(t, PasswordResetDataForTest.PasswordResetURLWithToken)

	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, emailVerifyLink)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}