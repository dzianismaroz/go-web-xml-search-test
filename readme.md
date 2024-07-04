Homework #4

Topic: XML data search and test coverage

This is a combined task on how to send requests, receive responses, work with parameters, headers, and also write tests.

The task is not difficult, the bulk of the work is writing various conditions and tests to satisfy these conditions.

We have some kind of search service:
* SearchClient - a structure with a FindUsers method that sends a request to an external system and returns the result, transforming it a little. It is located in the client.go file and cannot be edited.
* SearchServer is a kind of external system. Directly searches for data in the `dataset.xml` file. In production it would be launched as a separate web service, but in your code it will be launched as a separate handler.

Required:
* Write the SearchServer function in the file `client_test.go`, which you will run in the test through the test server (`httptest.NewServer`, example of use in `4/http/server_test.go`)
* Cover the FindUsers method with tests so that the coverage of the `client.go` file is as high as possible, namely 100%. Write tests in `client_test.go`. But when you run tests with the coverage flag, the total percentage will be written there, what percentage is in `client.go` - look in the report.
* It is also required to generate an HTML report with coverage. See an example of test coverage and report construction in `3/testing/coverage_test.go`.
* Tests must be written as complete ones, i.e. not to get coverage, but which actually test your code, check the returned result, edge cases, etc. They should show that SearchServer is working correctly.
* It follows from the previous paragraph that SearchServer also needs to be written as a full-fledged one

SearchServer accepts GET parameters:
* `query` - what to look for. We search in the `Name` and `About` record fields for just a substring, without regular characters. `Name` is first_name + last_name from xml (you need to manually go through the records in a loop and do this, you can’t do it automatically). If the field is empty, then we return all records (searching for an empty substring always returns true), i.e. we only do the sorting logic
* `order_field` - which field to sort by. It works by the fields `Id`, `Age`, `Name`, if empty, then we sort by `Name`, if something else, SearchServer complains with an error.
* `order_by` - sorting direction (as is, descending, ascending), client.go has corresponding constants
* `limit` - how many records to return
* `offset` - starting from which record to return (how much to skip from the beginning) - needed to organize page navigation

Additionally:
* Data for work is in the file `dataset.xml`
* How to work with XML - almost the same as with JSON, see the doc https://golang.org/pkg/encoding/xml/ and an example in the bot
* Run as `go test -cover`
* You can start by simply writing a server in `main.go` that implements the logic, and then move it to `client_test.go`
* Building coverage: `go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html`
* Documentation https://golang.org/pkg/net/http/ may help
* Use table testing - this is when you have a slice of test cases that differ in parameters.
* You may not be limited to the SearchServer function when testing if you need to check some completely separate tricky case, such as an error. But such cases will be few. Basically everything will be in SearchServer
* To cover one of the errors with a test, you will have to look into the source code of the function that returns this error and see under what operating conditions or input data this occurs. This is a client error, i.e. In this case, the request will not go to the server.
* Do not try to implement a timeout by connecting to an unknown IP
* The NextPage block on line 121 in client.go is used to create page navigation - I look in the server for the +1 entry - if there is one - I can show the next page

Code volume:
* SearchServer with all the structures and everything will be 170-200 lines
* Tests 200-300-400 lines, depending on the form - the main thing there will be a list of test cases

Recommended work plan:
1. Write code in the main function that simply implements the SearchServer logic based on fixed parameters and outputs it to the console, without http
2. Now format your code in an http handler, the parameters are no longer hardcoded, but taken from the request
3. Check with requests from the browser that the code works
4. Now start writing tests in client_test.go
5. First implement one test that simply makes a request through SearchClient to your HTTP handler running through the test server
6. Now build a report and see which code was called and which was not
7. Start writing test cases
8. Implement a separate handler or handlers for errors
