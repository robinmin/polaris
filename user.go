// UserModel defined a user model. It's an interface to the end user to enable them to extend the capability.
//
package polaris

import (
	"github.com/martini-contrib/sessionauth"
)

// UserModel can be any struct that represents a user in my system
type UserModel struct {
	Uid        int64  `form:"Uid" db:"N_UID"`
	RoleID     int    `form:"-" db:"N_ROLE_ID"`
	Grp01      int    `form:"-" db:"N_GRP1"`
	Grp02      int    `form:"-" db:"N_GRP2"`
	InUse      int    `form:"-" db:"N_INUSE"`
	UserID     string `form:"UserID" db:"C_UID"`
	Password   string `form:"Password" db:"C_PWD"`
	Email      string `form:"-" db:"C_EMAIL"`
	HomeURL    string `form:"-" db:"C_HOME"`
	Language   string `form:"-" db:"C_LANGUAGE"`
	Org01      string `form:"-" db:"C_ORG1"`
	Org02      string `form:"-" db:"C_ORG2"`
	Org03      string `form:"-" db:"C_ORG3"`
	Org04      string `form:"-" db:"C_ORG4"`
	Preference string `form:"-" db:"C_PREFERENCE"`

	authenticated bool `form:"-" db:"-"`
}

// GetAnonymousUser should generate an anonymous user model
// for all sessions. This should be an unauthenticated 0 value struct.
func GenerateAnonymousUser() sessionauth.User {
	return &UserModel{}
}

// Login will preform any actions that are required to make a user model
// officially authenticated.
func (u *UserModel) Login() {
	// Update last login time
	// Add to logged-in user's list
	// etc ...
	u.authenticated = true
}

// Logout will preform any actions that are required to completely
// logout a user.
func (u *UserModel) Logout() {
	// Remove from logged-in user's list
	// etc ...
	u.authenticated = false
}

func (u *UserModel) IsAuthenticated() bool {
	return u.authenticated
}

func (u *UserModel) UniqueId() interface{} {
	return u.Uid
}

// GetById will populate a user object from a database model with
// a matching id.
func (u *UserModel) GetById(id interface{}) error {
	// err := DbEngine.SelectOne(u, "select * from SYS_USER where N_UID = ?", id)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// GetByNamePass will populate a user object from a database model with
// a matching name and password.
func (u *UserModel) GetByNamePass(strUserName string, strPass string) error {
	// err := DbEngine.SelectOne(u, `select * from SYS_USER where C_UID = ? and C_PWD = substring(sys.fn_sqlvarbasetostr(hashbytes('MD5',?)),3,32)`, strUserName, strPass)
	// if err != nil {
	// 	return err
	// }
	// // correct user name & password have been inputed

	return nil
}
