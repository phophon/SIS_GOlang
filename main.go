package main
import (
   "fmt"
   "net/http"
   "encoding/json"
   _ "github.com/lib/pq"
   "database/sql"
   "log"
)

const (
   host     = "omega-postgresql-sgp1-08776-do-user-4090996-0.db.ondigitalocean.com"
   port     = 25061
   user     = "omega_rew"
   password = "c6eqgnwwv09cxlzo"
   dbname   = "TestPool"
   sslmode = "require"
 )

 type addressBook struct {
   Firstname string
   Lastname  string
   Code      int
   Phone     string
}

type EmergencyContact struct{
   Firstname string
   Lastname string
   Relationship string
   Phone string
   Email string
}

type Student struct {
   First_name  string
   Last_name      string
   Program string
   Cmkl_email string
   UUID int
   Photo string
   Contact struct {
      Phone_number string
      Personnal_email string
      Second_email string
   }
   Emergency []EmergencyContact
   Address struct{
      Addressstatus string
      City string
      State string
      Zip string
      Country string
   }
}

func main() {


   // fmt.Println(result)
   handleRequest()
}

func getAddressBookAll(w http.ResponseWriter, r *http.Request) {
   addBook := addressBook{
                Firstname: "Chaiyarin",
                Lastname:  "Niamsuwan",
                Code:      1993,
                Phone:     "0870940955",
              }
   json.NewEncoder(w).Encode(addBook)
}

func getAllStudent(w http.ResponseWriter, r *http.Request) {
   
   var uuid int
   var first_name string
   var last_name string
   var gender string
   var photo string
   var cmkl_email string
   var phone_number string
   var program string
   var personnal_email string
   var canvasid string
   var airtableid string
   var second_email string
   var studentList Student
   var address_id int
   var address string
   var city string
   var state string
   var zip string
   var country string

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

  result, err := db.Query("SELECT * FROM student")
   if err != nil {
      panic(err)
      log.Fatal(err)
      }

      for result.Next() {
         if err := result.Scan(&uuid, &first_name, &last_name, &gender, &photo, &phone_number, &cmkl_email, &program, &canvasid, &airtableid, &personnal_email, &second_email); err != nil {
            log.Fatal(err)
         }
      }

   resultE, err := db.Query("SELECT * FROM emergency WHERE uuid = 109877189")
      if err != nil {
         panic(err)
         log.Fatal(err)
         }
         
         var emergency_id int
         var first_nameE string
         var last_nameE string
         var relationship string
         var phone string
         var email string
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

         resultA, err := db.Query("SELECT * FROM address WHERE uuid = 109877189")
         if err != nil {
            panic(err)
            log.Fatal(err)
            }

            for resultA.Next() {
               if err := resultA.Scan(&address_id, &address, &city, &state, &zip, &country, &uuid); err != nil {
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

      fmt.Println(studentList)
      json.NewEncoder(w).Encode(studentList)
}

// func homePage(w http.ResponseWriter, r *http.Request) {
//    fmt.Fprint(w, "Welcome to the HomePage!")
// }

func handleRequest() {
   http.HandleFunc("/", getAllStudent)
   http.HandleFunc("/getAddress", getAddressBookAll)
   // http.HandleFunc("/getAllStudent", getAllStudent)
   http.ListenAndServe(":8910", nil)
}