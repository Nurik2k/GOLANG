document.getElementById("login-form").addEventListener("submit", function(event) {
    event.preventDefault();

    const login = document.getElementById("login").value;
    const password = document.getElementById("password").value;

    // Создаем объект с данными для отправки на сервер
    const data = {
        login: login,
      password: password
    };

    // Отправляем POST-запрос на сервер Golang
    fetch("/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
      alert(data.message); // Выводим сообщение от сервера
    })
    .catch(error => {
      console.error("Ошибка:", error);
    });
  });