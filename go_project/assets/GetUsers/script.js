// script.js

const baseUrl = "http://localhost:8080/users";
const deleteUrl = "http://localhost:8080/user";

function displayUsers() {
  fetch(baseUrl, {
    method: "GET",
  })
    .then((response) => response.json())
    .then((data) => {
      const table = document.getElementById("userTable");

      // Clear existing rows
      table.innerHTML = `
        <tr>
          <th>Login</th>
          <th>Password</th>
          <th>First Name</th>
          <th>Name</th>
          <th>Last Name</th>
          <th>Birthday</th>
          <th>Actions</th>
        </tr>
      `;

      // Display users
      data.forEach((user) => {
        const row = table.insertRow();

        const loginCell = row.insertCell();
        loginCell.textContent = user.login;

        const passwordCell = row.insertCell();
        passwordCell.textContent = user.password;

        const firstNameCell = row.insertCell();
        firstNameCell.textContent = user.first_name;

        const nameCell = row.insertCell();
        nameCell.textContent = user.name;

        const lastNameCell = row.insertCell();
        lastNameCell.textContent = user.last_name;

        const birthdayCell = row.insertCell();
        birthdayCell.textContent = user.birthday;

        const actionsCell = row.insertCell();
        const deleteButton = document.createElement("button");
        deleteButton.textContent = "Delete";
        deleteButton.addEventListener("click", () => deleteUser(user.id));
        actionsCell.appendChild(deleteButton);

        const updateButton = document.createElement("button");
        updateButton.textContent = "Edit";
        updateButton.addEventListener("click", () => {
          window.location.href = `file:///s%3A/GOLANG/test/Post/assets/PutUser/edit.html?id=${user.id}`;
        });
        actionsCell.appendChild(updateButton);
      });
    })
    .catch((error) => console.error("Error fetching users:", error));
}

// ... (deleteUser and updateUser functions remain unchanged)
function deleteUser(userId) {
  if (confirm("Are you sure you want to delete this user?")) {
    fetch(`${deleteUrl}/${userId}`, {
      method: "DELETE", 
    })
      .then(() => displayUsers())
      .catch((error) => console.error("Error deleting user:", error));
  }
}

// Display initial user table
displayUsers();
