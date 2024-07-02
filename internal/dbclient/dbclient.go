package dbclient

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"reflect"

	config "dealer-status/internal/config"

	"github.com/go-sql-driver/mysql"
)

type DbClient struct {
	User string
	Password string
	Host string
	Name string	
}

func ConnectToDatabase() driver.Connector{

	client := setCredentials()

	cfg := mysql.Config{
		User: client.User,
		Passwd: client.Password,
		Net: "tcp",
		Addr: client.Host,
		DBName: client.Name,
		AllowNativePasswords: true,
	}

	connector, connErr := mysql.NewConnector(&cfg)

	if connErr != nil {
		fmt.Println("Error Connecting to the Dealer Management Database: ", connErr.Error())
		os.Exit(1)
	}

	return connector
}

func OpenDatabaseAndTestConnection(conn driver.Connector) *sql.DB {
	db := sql.OpenDB(conn)

	if dbErr := db.Ping(); dbErr != nil {
		fmt.Println("Could not Ping Database: ", dbErr.Error())
		os.Exit(1)
	}

	fmt.Println("Database connection is successful.")

	return db
}

func setCredentials() DbClient {
	dbc := &DbClient{}
	
	attrs := reflect.TypeOf(*dbc)

	for i := 0 ; i < attrs.NumField() ; i++ {
		attrName := attrs.Field(i).Name
		attrVal := config.ValidateEnvCred(attrName, "DB")
		dbc.setProperty(attrName, attrVal)
	}

	return *dbc

}

func (c *DbClient) setProperty(name string, value string) {
	reflect.ValueOf(c).Elem().FieldByName(name).Set(reflect.ValueOf(value))
}

func GetWpQuery() string{
	db_prefix := config.ValidateEnvCred("PREFIX", "DB")

	return fmt.Sprintf(`
		select ID, pm1.meta_value, pm3.meta_value
		from %[1]vposts p
		join %[1]vpostmeta pm1 on p.ID = pm1.post_id
		join %[1]vpostmeta pm2 on p.ID = pm2.post_id
		join %[1]vpostmeta pm3 on p.ID = pm3.post_id
		where p.post_status = 'publish'
		and p.post_type = 'dealers'
		and pm1.meta_key = '_meta_url'
		and (
			pm2.meta_key = '_meta_is_live'
			and pm2.meta_value = 'true'
		)
		and pm3.meta_key = '_meta_ccid'
	`,
	db_prefix)
}