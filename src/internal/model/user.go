package model

import (
    "fmt"
)

type User struct {
    Email    string `json:"email"`
    FullName string `json:"fullname"`
    Password string `json:"password"`
}

var db database

func init() {
    db = New()
}

func (u *User) UpdateSessionData(sessionToken string, csrfToken string) error {
    _, err := db.Exec(
        "update users set session_token=?, csrf_token=? where email=?",
        sessionToken,
        csrfToken,
        u.Email,
    )
    if err != nil {
        return fmt.Errorf("AddSessionData: %v", err)
    }
    return nil
}

func (u *User) GetSessionData() (string, string, error) {
    var sessionToken string
    var csrfToken string
    err := db.QueryRow(
        "select session_token, csrf_token from users where email=?", u.Email,
        ).Scan(&sessionToken, &csrfToken)
    if err != nil {
        return "", "", fmt.Errorf("GetSessionData: %v", err)
    }
    return sessionToken, csrfToken, nil
}

func AddUser(user *User) error {
    _, err := db.Exec(
        "insert into users(email, fullname, password) values(?, ?, ?)",
        user.Email,
        user.FullName,
        user.Password,
        )
    if err != nil {
        return fmt.Errorf("AddUser: %v", err)
    }
    return nil
}

func DeleteUser(email string) error {
    _, err := db.Exec("delete from users where email=?", email)
    if err != nil {
        return fmt.Errorf("DeleteUser: %v", err)
    }
    return nil
}

func GetUsers() ([]User, error) {
    var users []User
    rows, err := db.Query("select email, fullname, password from users")
    if err != nil {
        return nil, fmt.Errorf("GetUsers: %v", err)
    }
    defer rows.Close()
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.Email, &user.FullName, &user.Password); err != nil {
            return nil, fmt.Errorf("GetUsers: %v", err)
        }
        users = append(users, user)
    }
    if err := rows.Err(); err != nil{
        return nil, fmt.Errorf("GetUsers: %v", err)
    }
    return users, nil
}

func GetUserById(id int) (User, error) {
    var user User
    err := db.QueryRow("select email, fullname, password from users where id=?", id).Scan(
        &user.Email,
        &user.FullName,
        &user.Password,
        )
    if err != nil {
        return User{}, fmt.Errorf("GetUserById: %v", err)
    }
    return user, nil
}

func GetUserByEmail(email string) (User, error) {
    var user User
    err := db.QueryRow("select email, fullname, password from users where email=?", email).Scan(
        &user.Email,
        &user.FullName,
        &user.Password,
        )
    if err != nil {
        return User{}, fmt.Errorf("GetUserByFullName: %v", err)
    }
    return user, nil
}

func UserExists(email string) bool {
    if _, err := GetUserByEmail(email); err != nil {
        return false
    }
    return true
}
