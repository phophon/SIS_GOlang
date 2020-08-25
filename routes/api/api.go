package api
import(
	"fmt"
	"net/http"
	"encoding/json"
	"database/sql"
	"log"
	"io/ioutil"
	"bytes"
	"strings"
	
	_ "github.com/lib/pq"
	"github.com/dgrijalva/jwt-go"

)

const (
	host     = "omega-postgresql-sgp1-08776-do-user-4090996-0.db.ondigitalocean.com"
	port     = 25061
	user     = "omega_rew"
	password = "c6eqgnwwv09cxlzo"
	dbname   = "TestPool"
	sslmode = "require"
  )
 
 type EmergencyContact struct {
	Firstname *string
	Lastname *string
	Relationship *string
	Phone *string
	Email *string
 }
 
 type Student struct {
	First_name  *string
	Last_name      *string
	Program *string
	Cmkl_email *string
	UUID int
	Photo *string
	Contact struct {
	   Phone_number *string
	   Personnal_email *string
	   Second_email *string
	}
	Emergency []EmergencyContact
	Address struct{
	   Addressstatus *string
	   City *string
	   State *string
	   Zip *string
	   Country *string
	}
 }

//  type Term struct{
// 	id int
// 	term_name string
// 	program string
//  }

 type Course struct{
	Id *string
	Course_name *string
	Schedule *string
	Unit int
	Room *string
	Instructor *string
	Status *string
 }

//  type EnrollStatus struct{
// 	status string
// 	message string
//  }

 type Enrollment struct{
	// term []Term
	Course []Course
	// enrollstatus []EnrollStatus
 } 

var ProfileApiHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var uuid int
	var first_name *string
	var last_name *string
	var gender *string
	var photo *string
	var cmkl_email *string
	var phone_number *string
	var program *string
	var personnal_email *string
	var canvasid *string
	var airtableid *string
	var second_email *string
	var studentList Student
	var address_id *int
	var address *string
	var city *string
	var state *string
	var zip *string
	var country *string
	var programid *int
	// id := 109877189
	ua := r.Header.Get("Authorization")
	fmt.Println("")
	fmt.Println("cilenttoken : ", ua)
	fmt.Println("")


	token, err := jwt.Parse(strings.Split(ua, " ")[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("my_secret_key"), nil
	})


	claims := token.Claims.(jwt.MapClaims);
	fmt.Println("===== claims :", claims["https://omega.auth/email"].(string))
	fmt.Println("claims passed")
 
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
 
   result, err := db.Query(`SELECT * FROM student WHERE cmkl_email = $1;`, claims["https://omega.auth/email"].(string))
	if err != nil {
	   panic(err)
	   log.Fatal(err)
	   }
 
	   for result.Next() {
		  if err := result.Scan(&uuid, &first_name, &last_name, &gender, &photo, &phone_number, &cmkl_email, &canvasid, &airtableid, &personnal_email, &second_email); err != nil {
			 log.Fatal(err)
		  }
	   }
 
	resultE, err := db.Query(`SELECT * FROM emergency WHERE uuid = $1;`, uuid)
	   if err != nil {
		  panic(err)
		  log.Fatal(err)
		  }
		  
		  var emergency_id *int
		  var first_nameE *string
		  var last_nameE *string
		  var relationship *string
		  var phone *string
		  var email *string
		  var emergencyContact EmergencyContact
 
		  for resultE.Next() {
			 if err := resultE.Scan(&emergency_id, &first_nameE, &last_nameE, &relationship, &phone, &email, &uuid); err != nil {
				log.Fatal(err)
			 }
			 emergencyContact.Firstname = first_nameE
			 emergencyContact.Lastname = last_nameE
			 emergencyContact.Relationship = relationship
			 emergencyContact.Phone = phone
			 emergencyContact.Email = email
			 studentList.Emergency = append(studentList.Emergency, emergencyContact)
		  }
 
		resultA, err := db.Query(`SELECT * FROM address WHERE uuid = $1;`, uuid)
		  if err != nil {
			 panic(err)
			 log.Fatal(err)
			 }
 
			 for resultA.Next() {
				if err := resultA.Scan(&address_id, &address, &city, &state, &zip, &country, &uuid); err != nil {
				   log.Fatal(err)
				}
			 }

		resultPE, err := db.Query(`SELECT * FROM programenrollment WHERE uuid = $1;`, uuid)
		  if err != nil {
			 panic(err)
			 log.Fatal(err)
			 }

			 var invoiceurl *string
			 var programenrollmentid *int
			 var registeredcredits *string
			 var status *bool
			 var type_ *string

			 for resultPE.Next() {
				if err := resultPE.Scan(&invoiceurl, &programenrollmentid, &registeredcredits, &status, &type_, &uuid, &programid); err != nil {
				   log.Fatal(err)
				}
			 }

		resultP, err := db.Query(`SELECT * FROM program WHERE programid = $1;`, programid)
		  if err != nil {
			 panic(err)
			 log.Fatal(err)
			 }

			 var shortname *string

			 for resultP.Next() {
				if err := resultP.Scan(&programid, &program, &airtableid, &shortname); err != nil {
				   log.Fatal(err)
				}
			 }

 
	   studentList.First_name = first_name
	   studentList.Last_name = last_name
	   studentList.Program = program
	   studentList.Cmkl_email = cmkl_email
	   studentList.UUID = uuid
	   studentList.Photo = photo
	   studentList.Contact.Phone_number = phone_number
	   studentList.Contact.Personnal_email = personnal_email
	   studentList.Contact.Second_email = second_email
	   studentList.Address.Addressstatus = address
	   studentList.Address.City = city
	   studentList.Address.State = state
	   studentList.Address.Zip = zip
	   studentList.Address.Country = country
 
	   w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(studentList)
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

type Profile struct{
	Name string
	Uuid int
}

var UpdateProfileHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var data Student
	var uuid int
	var address_uuid int
	var emergency_uuid int

	reqBody, err := json.Marshal(map[string]string{})

    resp, err := http.Post("http://localhost:8910/home",
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
    }
    
    defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
    // fmt.Println(string(body))
    fmt.Fprint(w, string(body))
	
	json.Unmarshal([]byte(string(body)), &data)
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
		   

		   if uuid == 0 {
			var programenrollmentid int
			var programid int

			result, err := db.Query(`SELECT programenrollmentid FROM programenrollment ORDER BY programenrollmentid DESC LIMIT 1;`)
			if err != nil {
			   panic(err)
			   log.Fatal(err)
			   }
			   for result.Next() {
				  if err := result.Scan(&programenrollmentid); err != nil {
					 log.Fatal(err)
				  }
			   }

			resultA, err := db.Query(`SELECT programid FROM program WHERE shortname = $1;`, data.Program)
			   if err != nil {
				  panic(err)
				  log.Fatal(err)
				}

				  for resultA.Next() {
					 if err := resultA.Scan(&programid); err != nil {
						log.Fatal(err)
					 }
				  }

			_, err = db.Exec(`INSERT INTO programenrollment (programenrollmentid, status, uuid, programid) values($1, $2, $3, $4);`, programenrollmentid+1, 1, data.UUID, programid)
				if err != nil {
					panic(err)
				}
				fmt.Println("inserted programenrollment")
		   }else{
			fmt.Println("updated programenrollment")
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

		if address_uuid == 0 {
			var address_id int
			result, err := db.Query(`SELECT address_id FROM address ORDER BY address_id DESC LIMIT 1;`)
				if err != nil {
				panic(err)
				log.Fatal(err)
				}
				for result.Next() {
					if err := result.Scan(&address_id); err != nil {
						log.Fatal(err)
					}
				}

			_, err = db.Exec(`INSERT INTO address (address_id, address, city, state, zip, country, uuid) values($1, $2, $3, $4, $5, $6, $7);`, address_id+1, data.Address.Addressstatus, data.Address.City, data.Address.State, data.Address.Zip, data.Address.Country, data.UUID)
				if err != nil {
					panic(err)
				}
				fmt.Println("inserted address")
			} else {
				sqlStatement := `UPDATE address SET address = $1, city = $2, state = $3, zip = $4, country = $5, uuid = $6 WHERE uuid = $7;`

				_, err = db.Exec(sqlStatement, data.Address.Addressstatus, data.Address.City, data.Address.State, data.Address.Zip, data.Address.Country, data.UUID, data.UUID)
					if err != nil {
						panic(err)
				}
				fmt.Println("updated address")
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

		   if emergency_uuid == 0 {
			   var emergency_id int

				result, err := db.Query(`SELECT emergency_id FROM emergency ORDER BY emergency_id DESC LIMIT 1;`)
				if err != nil {
				panic(err)
				log.Fatal(err)
				}
				for result.Next() {
					if err := result.Scan(&emergency_id); err != nil {
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
					fmt.Println("inserted address")
		   } else {
			var emergency_id int

			result, err := db.Query(`SELECT emergency_id FROM emergency WHERE uuid = $1;`, data.UUID)
			if err != nil {
			panic(err)
			log.Fatal(err)
			}
			for result.Next() {
				if err := result.Scan(&emergency_id); err != nil {
					log.Fatal(err)
				}
			}

			sqlStatement := `UPDATE emergency SET first_name = $1, last_name = $2, relationship = $3, phone = $4, email = $5, uuid = $6 WHERE emergency_id = $7;`
			for i, s := range data.Emergency {
				_, err = db.Exec(sqlStatement, s.Firstname, s.Lastname, s.Relationship, s.Phone, s.Email, data.UUID, emergency_id-1+i)
				if err != nil {
					panic(err)
				}
			}
			fmt.Println("updated address")

		   }
})