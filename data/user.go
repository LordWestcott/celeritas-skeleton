package data

import (
	"errors"
	"time"

	"github.com/lordwestcott/celeritas"
	up "github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `db:"id,omitempty"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Active    int       `db:"user_active"`
	Password  string    `db:"password"`
	Token     Token     `db:"-"` //Doesn't exist in the database
}

func (u *User) Validate(validator *celeritas.Validation) {
	validator.Check(u.LastName != "", "last_name", "Last Name is required")
	validator.Check(u.FirstName != "", "first_name", "First Name is required")
	validator.Check(u.Email != "", "email", "Email is required")
	validator.IsEmail("email", u.Email)
}

// Override the table name
func (u *User) Table() string {
	return "users"
}

func (u *User) GetAll() ([]*User, error) {
	collection := upper.Collection(u.Table())

	var all []*User

	res := collection.Find().OrderBy("last_name")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	var user User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"email =": email})
	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id =": user.ID, "expiry >": time.Now()}).OrderBy("created_at DESC")
	err = res.One(&token)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}
	user.Token = token

	return &user, nil
}

func (u *User) Get(id int) (*User, error) {
	var user User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id =": user.ID, "expiry >": time.Now()}).OrderBy("created_at DESC")
	err = res.One(&token)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}
	user.Token = token

	return &user, nil
}

func (u *User) Update(user User) error {
	user.UpdatedAt = time.Now()
	collection := upper.Collection(u.Table())
	res := collection.Find(user.ID)
	err := res.Update(&user)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(id int) error {
	collection := upper.Collection(u.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Insert(user User) (int, error) {
	newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Password = string(newHash)
	collection := upper.Collection(u.Table())
	res, err := collection.Insert(&user)
	if err != nil {
		return 0, err
	}

	return getInsertID(res.ID()), nil
}

func (u *User) ResetPassword(id int, password string) error {
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user, err := u.Get(id)
	if err != nil {
		return err
	}

	u.Password = string(newHash)
	err = user.Update(*u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (u *User) CheckForRememberToken(id int, token string) bool {
	var rememberToken RememberToken
	rt := RememberToken{}
	collection := upper.Collection(rt.Table())
	cond := up.Cond{"user_id =": id, "remember_token =": token}
	res := collection.Find(cond)
	err := res.One(&rememberToken)
	return err == nil
}
