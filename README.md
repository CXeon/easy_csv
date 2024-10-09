easy_csv
====

The easy_csv package is used to marshal or unmarshal csv file data.  

Installation
===

```go get -u github.com/CXeon/easy_csv@v1.0.1```  

Full example
===

Marshal a structure to a csv file.  
---

```golang
package main

import (
	"github.com/CXeon/easy_csv"
	"os"
)

type testStudentInfo struct {
	Name  string `csv:"name"`
	Age   int
	Grade string
	Score float64 `csv:"mScore"`
	Email string  `csv:"mEmail,email_desensitization"` //The email_desensitization will desensitize email addresses.The email username must be greater than 3 characters to be effective.
	Phone string  `csv:"mPhone,phone_desensitization"` //The phone_desensitization will desensitize the phone number.The mobile phone number must be greater than 6 digits to be effective.
}

func main() {
	fileName := "./CsvWriter.csv"

	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	clientWriter := easy_csv.NewClientWriter(csvFile)

	//clientWriter.WriteString2File()

	row := testStudentInfo{
		Name:  "Jam",
		Age:   10,
		Grade: "Grade 4",
		Score: 99.01,
		Email: "jamjam@test.com",
		Phone: "13322226666",
	}

	// parameter row cloud be &row
	err = clientWriter.WriteRow2File(row, true)
	if err != nil {
		panic(err)
	}

}

```
The result is  

|name|Age|Grade|mScore|mEmail|mPhone|
|---|---|---|---|---|---|
|Jam|10|Grade 4|99.01|ja***m@test.com|133****6666|

Marshal a list to a csv file.
---

```golang
package main

import (
	"github.com/CXeon/easy_csv"
	"os"
)

type testStudentInfo struct {
	Name  string `csv:"name"`
	Age   int
	Grade string
	Score float64 `csv:"mScore"`
	Email string  `csv:"mEmail,email_desensitization"` //The email_desensitization will desensitize email addresses.The email username must be greater than 3 characters to be effective.
	Phone string  `csv:"mPhone,phone_desensitization"` //The phone_desensitization will desensitize the phone number.The mobile phone number must be greater than 6 digits to be effective.
}

func main() {
	fileName := "./CsvWriter.csv"

	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	clientWriter := easy_csv.NewClientWriter(csvFile)

	list := []testStudentInfo{
		{
			Name:  "Jam",
			Age:   10,
			Grade: "Grade 4",
			Score: 99.01,
			Email: "jamjam@test.com",
			Phone: "13322226666",
		},
		{
			Name:  "xeon",
			Age:   13,
			Grade: "Grade 6",
			Score: 101.11,
			Email: "xeonjoe@test.com",
			Phone: "16542354654",
		},
	}

	// parameter list  cloud be &list.the item of list cloud be a structure pointer
	err = clientWriter.WriteRows2File(list, true)
	if err != nil {
		panic(err)
	}

}

```

The result is  

|name|Age|Grade|mScore|mEmail|mPhone|
|---|---|---|---|---|---|
|Jam|10|Grade 4|99.01|ja***m@test.com|133****6666|
|xeon|13|Grade 6|101.11|xe***e@test.com|165****4654|

Unmarshal a row of a csv to a structure
---

a csv is 

| name |Age|Grade|mScore| mEmail |mPhone|
|------|---|---|---|---------------|---|
| Jam  |10|Grade 4|99.01| ja***m@test.com |133****6666|
| xeon |13|Grade 6|101.11| xe***e@test.com |165****4654|
| bob  |10|Grade 4|99.01| bo***i@test.com |133****6666|

```golang
package main

import (
	"fmt"
	"github.com/CXeon/easy_csv"
	"os"
)

type testStudentInfo struct {
	Name  string 
	Age   int
	Grade string
	Score float64 
	Email string   
	Phone string  
}

func main() {
	fileName := "./CsvReader.csv"

	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	clientReader := easy_csv.NewClientReader(csvFile)

	clientReader.Read() //Important.Exclude title

	data := testStudentInfo{}
	err = clientReader.ReadRowFromFile(&data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", data)
}

```
the result show  

main.testStudentInfo{Name:"Jam", Age:10, Grade:"Grade 4", Score:99.01, Email:"ja***m@test.com", Phone:"133****6666"}


If the order of your csv file columns and structure fields is inconsistent
---

```golang
package main

import (
	"fmt"
	"github.com/CXeon/easy_csv"
	"os"
)

type testStudentInfo struct {
	Score float64 
	Email string 
	Phone string  
	Name  string 
	Age   int
	Grade string
}

func main() {
	fileName := "./CsvReader.csv"

	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	clientReader := easy_csv.NewClientReader(csvFile)

	clientReader.Read() //Important.Exclude title

	data := testStudentInfo{}

	names := []string{
		"Name",
		"Age",
		"Grade",
		"Score",
		"Email",
		"Phone",
	}
	err = clientReader.ReadRowFromFileWithNames(names, &data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", data)
}

```

The result show

main.testStudentInfo{Score:99.01, Email:"ja***m@test.com", Phone:"133****6666", Name:"Jam", Age:10, Grade:"Grade 4"}


Unmarshal rows of a csv to a structure
---

a csv is

| name |Age|Grade|mScore| mEmail|mPhone|
|------|---|---|---|---------------|---|
| Jam  |10|Grade 4|99.01| ja***m@test.com |133****6666|
| xeon |13|Grade 6|101.11| xe***e@test.com |165****4654|
| bob  |10|Grade 4|99.01| bo***i@test.com |133****6666|

```golang
package main

import (
	"fmt"
	"github.com/CXeon/easy_csv"
	"os"
)

type testStudentInfo struct {
	Name  string 
	Age   int
	Grade string
	Score float64 
	Email string  
	Phone string  
}

func main() {
	fileName := "./CsvReader.csv"

	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	clientReader := easy_csv.NewClientReader(csvFile)

	clientReader.Read() //Important.Exclude title

	var list []testStudentInfo
	// var list  []*testStudentInfo //yes
	err = clientReader.ReadRowsFromFile(&list)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", list)
}

```

The result show

[]main.testStudentInfo{main.testStudentInfo{Name:"Jam", Age:10, Grade:"Grade 4", Score:99.01, Email:"ja***m@test.com", Phone:"133****6666"}, main.testStudentInfo{Name:"xeon", Age:13, Grade:"Grade 5", Score:101.11, Email:"xe***e@test.com", Phone:"165****4654"}, main.testStudentInfo{Name:"bob", Age:10, Grade:"Grade 4", Score:99.01, Email:"bo***i@test.com", Phone:"133****6666"}}  

If the order of your csv file columns and structure fields is inconsistent
---

```golang
package main

import (
	"fmt"
	"github.com/CXeon/easy_csv"
	"os"
)

type testStudentInfo struct {
	Score float64
	Email string  
	Phone string  
	Name  string 
	Age   int
	Grade string
}

func main() {
	fileName := "./CsvReader.csv"

	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	clientReader := easy_csv.NewClientReader(csvFile)

	clientReader.Read() //Important.Exclude title

	var list []testStudentInfo
	// var list  []*testStudentInfo //yes
	names := []string{
		"Name",
		"Age",
		"Grade",
		"Score",
		"Email",
		"Phone",
	}
	err = clientReader.ReadRowsFromFileWithNames(names, &list)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", list)
}

```

The result show  

[]main.testStudentInfo{main.testStudentInfo{Score:99.01, Email:"ja***m@test.com", Phone:"133****6666", Name:"Jam", Age:10, Grade:"Grade 4"}, main.testStudentInfo{Score:101.11, Email:"xe***e@test.com", Phone:"165****4654", Name:"xeon", Age:13, Grade:"Grade 5"}, main.testStudentInfo{Score:99.01, Email:"bo***i@test.com", Phone:"133****6666", Name:"bob", Age:10, Grade:"Grade 4"}}
