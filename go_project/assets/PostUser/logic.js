document.getElementById("AddUser-form").addEventListener("submit", function (event) {
    event.preventDefault();

    const login = document.getElementById("login").value;
    const password = document.getElementById("password").value;
    const first_name = document.getElementById("first_name").value
    const name = document.getElementById("name").value
    const last_name = document.getElementById("last_name").value
    const birthday = document.getElementById("birthday").value

    // Создаем объект с данными для отправки на сервер
    const data = {
        login: login,
        password: password,
        first_name: first_name,
        name: name,
        last_name: last_name,
        birthday: birthday
      };

    const url = "http://localhost:8080/user";

    const options = {
        method: "POST",
        body: JSON.stringify(data),
    };

  // Отправляем POST-запрос на сервер Golang
    fetch(url, options) 
      .then(response => {console.log(response)})
      .then(data=>{window.location.href = "file:///s%3A/GOLANG/test/Post/assets/GetUsers/UsersList.html"})
      .catch(error => {
        console.log("Ошибка:", error);
      });
});

