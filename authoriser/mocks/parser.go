// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"github.com/ONSdigital/dp-authorisation/v2/permissions"
	"sync"
)

// TokenParserMock is a mock implementation of authoriser.TokenParser.
//
// 	func TestSomethingThatUsesTokenParser(t *testing.T) {
//
// 		// make and configure a mocked authoriser.TokenParser
// 		mockedTokenParser := &TokenParserMock{
// 			ParseFunc: func(tokenString string) (*permissions.EntityData, error) {
// 				panic("mock out the Parse method")
// 			},
// 		}
//
// 		// use mockedTokenParser in code that requires authoriser.TokenParser
// 		// and then make assertions.
//
// 	}
type TokenParserMock struct {
	// ParseFunc mocks the Parse method.
	ParseFunc func(tokenString string) (*permissions.EntityData, error)

	// calls tracks calls to the methods.
	calls struct {
		// Parse holds details about calls to the Parse method.
		Parse []struct {
			// TokenString is the tokenString argument value.
			TokenString string
		}
	}
	lockParse sync.RWMutex
}

// Parse calls ParseFunc.
func (mock *TokenParserMock) Parse(tokenString string) (*permissions.EntityData, error) {
	if mock.ParseFunc == nil {
		panic("TokenParserMock.ParseFunc: method is nil but TokenParser.Parse was just called")
	}
	callInfo := struct {
		TokenString string
	}{
		TokenString: tokenString,
	}
	mock.lockParse.Lock()
	mock.calls.Parse = append(mock.calls.Parse, callInfo)
	mock.lockParse.Unlock()
	return mock.ParseFunc(tokenString)
}

// ParseCalls gets all the calls that were made to Parse.
// Check the length with:
//     len(mockedTokenParser.ParseCalls())
func (mock *TokenParserMock) ParseCalls() []struct {
	TokenString string
} {
	var calls []struct {
		TokenString string
	}
	mock.lockParse.RLock()
	calls = mock.calls.Parse
	mock.lockParse.RUnlock()
	return calls
}