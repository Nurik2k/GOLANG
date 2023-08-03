document.getElementById("login-form").addEventListener("submit", function (event) {
    event.preventDefault();

    const login = document.getElementById("login").value;
    const password = document.getElementById("password").value;

    // Создаем объект с данными для отправки на сервер
    const data = {
        login: login,
        password: password,
      };

    const url = "http://localhost:8080/user";

    const options = {
        method: "POST",
    };

  // Отправляем POST-запрос на сервер Golang
    fetch(url, options) 
      .then(response => {console.log(response)})
      .then(data => window.location.href = "file:///s%3A/GOLANG/test/Post/assets/GetUsers/UsersList.html")
      .catch(error => {
        console.log("Ошибка:", error);
      });
});

