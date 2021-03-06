package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

const (
	host     = "omega-postgresql-sgp1-08776-do-user-4090996-0.db.ondigitalocean.com"
	port     = 25061
	user     = "omega_rew"
	password = "c6eqgnwwv09cxlzo"
	dbname   = "TestPool"
	sslmode  = "require"
)

type EmergencyContact struct {
	Firstname    string
	Lastname     string
	Relationship string
	Phone        string
	Email        string
}

type Address struct {
	Addressstatus string
	City          string
	State         string
	Zip           string
	Country       string
}

type Student struct {
	First_name string
	Last_name  string
	Program    string
	Cmkl_email string
	UUID       string
	Photo      string
	Contact    struct {
		Phone_number    string
		Personnal_email string
		Second_email    string
	}
	Emergency   [2]EmergencyContact
	Useraddress [3]Address
}

//  type Term struct{
// 	id int
// 	term_name string
// 	program string
//  }

type Course struct {
	Id          *string
	Course_name *string
	Schedule    *string
	Unit        int
	Room        *string
	Instructor  *string
	Status      *string
}

//  type EnrollStatus struct{
// 	status string
// 	message string
//  }

type Enrollment struct {
	// term []Term
	Course []Course
	// enrollstatus []EnrollStatus
}

var ProfileApiHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var uuid sql.NullString
	var studentid sql.NullString
	var first_name sql.NullString
	var last_name sql.NullString
	var gender sql.NullString
	var photo sql.NullString
	var cmkl_email sql.NullString
	var phone_number sql.NullString
	var program sql.NullString
	var personnal_email sql.NullString
	var canvasid sql.NullString
	var airtableid sql.NullString
	var second_email sql.NullString
	var studentList Student
	var address_id int
	var addressstatus sql.NullString
	var city sql.NullString
	var state sql.NullString
	var zip sql.NullString
	var country sql.NullString
	var roles sql.NullString
	var programid int
	// id := 109877189
	ua := r.Header.Get("Authorization")
	fmt.Println("")
	fmt.Println("cilenttoken : ", ua)
	fmt.Println("")

	token, err := jwt.Parse(strings.Split(ua, " ")[1], func(token *jwt.Token) (interface{}, error) {
		return []byte("secureSecretText"), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Couldn't parse claims")
	}
	fmt.Println(claims)

	// if claims["exp"].(int64) < time.Now().UTC().Unix() {
	// 	fmt.Println("JWT is expired")
	// }
	// fmt.Println("===== claims :", claims["https://omega.auth/email"].(string))
	// fmt.Println("claims passed")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	result, err := db.Query(`SELECT * FROM student WHERE cmkl_email = $1;`, claims["CmklMail"].(string))
	if err != nil {
		panic(err)
		log.Fatal(err)
	}

	for result.Next() {
		if err := result.Scan(&first_name, &last_name, &gender, &photo, &phone_number, &cmkl_email, &canvasid, &airtableid, &personnal_email, &second_email, &uuid, &roles, &studentid); err != nil {
			fmt.Println("=========Error=================")
			log.Fatal(err)
			fmt.Println("=========Error End=================")
		}
	}
	fmt.Println("uuid =====", uuid)
	resultE, err := db.Query(`SELECT * FROM emergency WHERE uuid = $1;`, uuid.String)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}

	var emergency_id int
	var count = 0
	var first_nameE sql.NullString
	var last_nameE sql.NullString
	var relationship sql.NullString
	var phone sql.NullString
	var email sql.NullString
	var emergencyContact EmergencyContact
	for resultE.Next() {
		if err := resultE.Scan(&emergency_id, &first_nameE, &last_nameE, &relationship, &phone, &email, &uuid); err != nil {
			log.Fatal(err)
		}
		emergencyContact.Firstname = first_nameE.String
		emergencyContact.Lastname = last_nameE.String
		emergencyContact.Relationship = relationship.String
		emergencyContact.Phone = phone.String
		emergencyContact.Email = email.String
		studentList.Emergency[count] = emergencyContact
		count += 1
	}

	fmt.Println("uuid =====", uuid)
	resultA, err := db.Query(`SELECT * FROM address WHERE uuid = $1;`, uuid.String)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}
	fmt.Println("resultSet =====", resultA)
	var countA = 0
	var address Address

	for resultA.Next() {
		if err := resultA.Scan(&address_id, &addressstatus, &city, &state, &zip, &country, &uuid); err != nil {
			log.Fatal(err)
		}
		address.Addressstatus = addressstatus.String
		address.City = city.String
		address.State = state.String
		address.Zip = zip.String
		address.Country = country.String
		studentList.Useraddress[0] = address
		studentList.Useraddress[1] = address
		studentList.Useraddress[2] = address
		countA++
	}

	resultPE, err := db.Query(`SELECT * FROM programenrollment WHERE uuid = $1;`, uuid)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}

	var invoiceurl *string
	var programenrollmentid int
	var registeredcredits *string
	var status bool
	var type_ *string

	for resultPE.Next() {
		if err := resultPE.Scan(&invoiceurl, &programenrollmentid, &registeredcredits, &status, &type_, &programid, &uuid); err != nil {
			log.Fatal(err)
		}
	}

	resultP, err := db.Query(`SELECT * FROM program WHERE programid = $1;`, programid)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}

	var shortname string

	for resultP.Next() {
		if err := resultP.Scan(&programid, &program, &airtableid, &shortname); err != nil {
			log.Fatal(err)
		}
	}

	studentList.First_name = first_name.String
	studentList.Last_name = last_name.String
	studentList.Program = program.String
	studentList.Cmkl_email = cmkl_email.String
	studentList.UUID = studentid.String
	studentList.Photo = photo.String
	studentList.Contact.Phone_number = phone_number.String
	studentList.Contact.Personnal_email = personnal_email.String
	studentList.Contact.Second_email = second_email.String
	//    studentList.Address.Addressstatus = address.String
	//    studentList.Address.City = city.String
	//    studentList.Address.State = state.String
	//    studentList.Address.Zip = zip.String
	//    studentList.Address.Country = country.String

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(studentList)
	fmt.Println(studentList)
})

var EnrollmentApiHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var enrollmentList Enrollment

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	//    result, err := db.Query("SELECT * FROM semester")
	//    if err != nil {
	// 	panic(err)
	// 	log.Fatal(err)
	// 	}

	// 	   for result.Next() {

	// 		var term Term
	// 		var semesterid int
	// 		var semestername string
	// 		var startdate string
	// 		var enddate string
	// 		var airtableid string
	// 		var academicyearid int

	// 		  if err := result.Scan(&semesterid, &semestername, &startdate, &enddate, &airtableid, &academicyearid); err != nil {
	// 			 log.Fatal(err)
	// 		  }
	// 		  term.id = semesterid
	// 		  term.term_name = semestername
	// 		  enrollmentList.term = append(enrollmentList.term, term)
	//        }

	resultC, err := db.Query("SELECT * FROM course")
	if err != nil {
		panic(err)
		log.Fatal(err)
	}

	for resultC.Next() {

		var courseid int
		var code *string
		var description *string
		var name *string
		var airtableid *string
		var unit int
		var room *string
		var status *string
		var time *string
		var instructor *string
		var course Course

		if err := resultC.Scan(&courseid, &code, &description, &name, &airtableid, &unit, &room, &status, &time, &instructor); err != nil {
			log.Fatal(err)
		}
		course.Id = code
		course.Course_name = name
		course.Schedule = time
		course.Unit = unit
		course.Room = room
		course.Status = status
		course.Instructor = instructor
		enrollmentList.Course = append(enrollmentList.Course, course)
	}

	// 	resultE, err := db.Query("SELECT * FROM courseenrollment")
	//    if err != nil {
	// 	panic(err)
	// 	log.Fatal(err)
	// 	}

	// 	for resultE.Next() {

	// 		var courseofferid int
	// 		var uuid string
	// 		var status string
	// 		var message string
	// 		var enrollstatus EnrollStatus

	// 		if err := resultC.Scan(&courseofferid, &uuid, &status, &message); err != nil {
	// 			log.Fatal(err)
	// 		}

	// 		enrollstatus.status = status
	// 		enrollstatus.message = message
	// 		enrollmentList.enrollstatus = append(enrollmentList.enrollstatus, enrollstatus)
	// 	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(enrollmentList)
})

var UpdateProfileHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var data Student
	var uuid sql.NullString
	var address_uuid sql.NullString
	var emergency_uuid sql.NullString

	json.NewDecoder(r.Body).Decode(&data)

	// reqBody, err := json.Marshal(map[string]string{})

	// resp, err := http.Post("http://localhost:8910/api/v1/home",
	// 	"application/json", bytes.NewBuffer(reqBody))
	// if err != nil {
	// 	print(err)
	// }

	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	print(err)
	// }
	fmt.Println(data)
	fmt.Fprint(w, data)

	// json.Unmarshal([]byte(string(body)), &data)
	// fmt.Println(data)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `UPDATE student SET uuid = $1, first_name = $2, last_name = $3, cmkl_email = $4, photo = $5, phone_number = $6, personnal_email = $7, second_email = $8 WHERE cmkl_email = $9;`

	_, err = db.Exec(sqlStatement, data.UUID, data.First_name, data.Last_name, data.Cmkl_email, data.Photo, data.Contact.Phone_number, data.Contact.Personnal_email, data.Contact.Second_email, data.Cmkl_email)
	if err != nil {
		panic(err)
	}
	fmt.Println("Update Profile")

	result, err := db.Query(`SELECT uuid FROM programenrollment WHERE uuid = $1;`, data.UUID)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}

	for result.Next() {
		if err := result.Scan(&uuid); err != nil {
			log.Fatal(err)
		}
	}

	if uuid.String == " " {
		var programenrollmentid int
		var programid int

		resultPE, err := db.Query(`SELECT programenrollmentid FROM programenrollment ORDER BY programenrollmentid DESC LIMIT 1;`)
		if err != nil {
			panic(err)
			log.Fatal(err)
		}
		for resultPE.Next() {
			if err := resultPE.Scan(&programenrollmentid); err != nil {
				log.Fatal(err)
			}
		}

		resultP, err := db.Query(`SELECT programid FROM program WHERE shortname = $1;`, data.Program)
		if err != nil {
			panic(err)
			log.Fatal(err)
		}

		for resultP.Next() {
			if err := resultP.Scan(&programid); err != nil {
				log.Fatal(err)
			}
		}

		_, err = db.Exec(`INSERT INTO programenrollment (programenrollmentid, status, uuid, programid) values($1, $2, $3, $4);`, programenrollmentid+1, 1, data.UUID, programid)
		if err != nil {
			panic(err)
		}
		fmt.Println("Inserted ProgramEnrollment")
	} else {
		fmt.Println("Updated ProgramEnrollment")
	}

	resultA, err := db.Query(`SELECT uuid FROM address WHERE uuid = $1;`, data.UUID)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}

	for resultA.Next() {
		if err := resultA.Scan(&address_uuid); err != nil {
			log.Fatal(err)
		}
	}

	if address_uuid.String == "" {
		var address_id int
		resultA2, err := db.Query(`SELECT address_id FROM address ORDER BY address_id DESC LIMIT 1;`)
		if err != nil {
			panic(err)
			log.Fatal(err)
		}
		for resultA2.Next() {
			if err := resultA2.Scan(&address_id); err != nil {
				log.Fatal(err)
			}
		}
		for i, s := range data.Useraddress {
			sqlStatement := `INSERT INTO address (address_id, address_status, city, state, zip, country, uuid) values($1, $2, $3, $4, $5, $6, $7);`
			_, err = db.Exec(sqlStatement, address_id+1+i, s.Addressstatus, s.City, s.State, s.Zip, s.Country, data.UUID)
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("Inserted Address")
	} else {
		var AddressID int
		var listAddressID []int
		resultA2, err := db.Query(`SELECT address_id FROM address;`)
		if err != nil {
			panic(err)
			log.Fatal(err)
		}
		for resultA2.Next() {
			if err := resultA2.Scan(&AddressID); err != nil {
				log.Fatal(err)
			} else {
				listAddressID = append(listAddressID, AddressID)
			}
		}

		sqlStatement := `UPDATE address SET address_status = $1, city = $2, state = $3, zip = $4, country = $5 WHERE address_id = $6;`
		for i, s := range data.Useraddress {
			_, err = db.Exec(sqlStatement, s.Addressstatus, s.City, s.State, s.Zip, s.Country, listAddressID[i])
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("Updated Address")
	}

	resultE, err := db.Query(`SELECT uuid FROM emergency WHERE uuid = $1;`, data.UUID)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}
	for resultE.Next() {
		if err := resultE.Scan(&emergency_uuid); err != nil {
			log.Fatal(err)
		}
	}

	if emergency_uuid.String == "" {
		var emergency_id int

		resultE2, err := db.Query(`SELECT emergency_id FROM emergency ORDER BY emergency_id DESC LIMIT 1;`)
		if err != nil {
			panic(err)
			log.Fatal(err)
		}
		for resultE2.Next() {
			if err := resultE2.Scan(&emergency_id); err != nil {
				log.Fatal(err)
			}
		}
		sqlStatement := `INSERT INTO emergency (emergency_id, first_name, last_name, relationship, phone, email, uuid) values($1, $2, $3, $4, $5, $6, $7);`
		for i, s := range data.Emergency {
			_, err = db.Exec(sqlStatement, emergency_id+1+i, s.Firstname, s.Lastname, s.Relationship, s.Phone, s.Email, data.UUID)
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("Inserted Emergency")
	} else {
		var EmergencyID int
		var ListEmergencyID []int

		resultE3, err := db.Query(`SELECT emergency_id FROM emergency WHERE uuid = $1;`, data.UUID)
		if err != nil {
			panic(err)
			log.Fatal(err)
		}
		for resultE3.Next() {
			if err := resultE3.Scan(&EmergencyID); err != nil {
				log.Fatal(err)
			}
			ListEmergencyID = append(ListEmergencyID, EmergencyID)
		}

		sqlStatement := `UPDATE emergency SET first_name = $1, last_name = $2, relationship = $3, phone = $4, email = $5, uuid = $6 WHERE emergency_id = $7;`
		for i, s := range data.Emergency {
			_, err = db.Exec(sqlStatement, s.Firstname, s.Lastname, s.Relationship, s.Phone, s.Email, data.UUID, ListEmergencyID[i])
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("Updated Emergency")
	}
})
