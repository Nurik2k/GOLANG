document.getElementById("login-form").addEventListener("submit", function(event) {
    event.preventDefault();

    const login = document.getElementById("login").value;
    const password = document.getElementById("password").value;

    // Создаем объект с данными для отправки на сервер
    const data = {
        login: login,
      password: password
    };

    const url = "http://localhost:8080/login";

    const options = {
      method: "POST",
      body: JSON.stringify(data),
      headers: {
        "Content-Type": "application/json"
      },
      mode: "no-cors" // Set 'no-cors' mode to bypass CORS restrictions"
    };

    // Отправляем POST-запрос на сервер Golang
    fetch(url, options)
    .then(response => response.json())
    .then(data => {
      alert(data.message); // Выводим сообщение от сервера
    })
    .catch(error => {
      console.log("Ошибка:", error);
    });
  });