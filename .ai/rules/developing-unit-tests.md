---
description: When writing unit tests
alwaysApply: false
---
- Always follow uber style guide for go.
- Use testify for assertions
- Use testify/mock for mocking
- Write table-driven tests
- Tests tables should follow the format of `map[string]struct{...}` where the string key is the name of the test and
   the struct contains the input and expected output of the test. 
- Tests names should be descriptive and follow the format of `expects/should ... when`. 
- Tests should follow the AAA pattern, and we should always indicate with comments where which section starts: `Arrange, Act, Assert`.
- Tests should be developed under the package `package <feature>_test` and should be placed in the same directory as the package they are testing. 
- You must NOT use conditional statements on tests 
- Use assert.ErrorIs to compare errors and don't compare if error is nil
- Use assert.AnError to mock errors.
- Create mocks inside mock_test.go file.
- Don't use mock.Anything 
- Before you create a new mock, check if it already exists
- Always use Address-of operator to create a pointer to the mock
- Group variables with var when is more than one, for example: group variables in Arrange
- Add expected values to test table.
- Separate test table struct in input, mocks and expected amd pre declare types for all. 
- Add AssertExpectations on mocks
- Use expected to name expected values, for example expectedOutput, expectedErr 
- Ignore ctx when call m.Called method on mocks