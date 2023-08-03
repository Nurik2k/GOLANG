const baseUrl = "http://localhost:8080/user";
const editUrl = "http://localhost:8080/user";

function getUserById(userId) {
  fetch(`${baseUrl}/${userId}`, {
    method: "GET",
  })
    .then((response) => response.json())
    .then((data) => fillEditForm(data))
    .catch((error) => console.error("Error fetching user by id:", error));
}

function fillEditForm(user) {
  const editForm = document.getElementById("editUserForm");
  editForm.login.value = user.login;
  editForm.password.value = user.password;
  editForm.first_name.value = user.first_name;
  editForm.name.value = user.name;
  editForm.last_name.value = user.last_name;
  editForm.birthday.value = user.birthday;

  editForm.onsubmit = function (event) {
    event.preventDefault();
    updateUser(user.id, editForm);
  };
}

function updateUser(userId, editForm) {
  const userData = {
    login: editForm.login.value,
    password: editForm.password.value,
    first_name: editForm.first_name.value,
    name: editForm.name.value,
    last_name: editForm.last_name.value,
    birthday: editForm.birthday.value,
  };

  fetch(`${editUrl}/${userId}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(userData),
  })
    .then(() => window.location.href = "file:///s%3A/GOLANG/test/Post/assets/GetUsers/UsersList.html")
    .catch((error) => console.error("Error updating user:", error));
}

// Get the user ID from the URL query parameter
const urlParams = new URLSearchParams(window.location.search);
const userId = urlParams.get("id");
if (userId) {
  getUserById(userId);
}
