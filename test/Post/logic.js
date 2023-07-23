document.getElementById("login-form").addEventListener("submit", function(event) {
    event.preventDefault();

    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    // Создаем объект с данными для отправки на сервер
    const data = {
        username: username,
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