document.addEventListener("DOMContentLoaded", function () {
  const userList = document.getElementById("userList");
  const refreshUserListButton = document.getElementById("refreshUserList");

  const createUserForm = document.getElementById("createUserForm");
  const createUserMessage = document.getElementById("createUserMessage");

  const updateUserSection = document.getElementById("updateUserSection");
  const updateUserForm = document.getElementById("updateUserForm");
  const updateUserMessage = document.getElementById("updateUserMessage");
  const cancelUpdateButton = document.getElementById("cancelUpdate");

  let authToken = localStorage.getItem("authToken");

  // Fetch Users
  async function fetchUsers() {
    userList.innerHTML = "Loading...";
    try {
      const response = await fetch("http://localhost:4000/v1/user", {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      renderUserList(data.users);
    } catch (error) {
      userList.innerHTML = "Error loading users";
      console.error("Error fetching users:", error);
    }
  }

  // Render User List
  function renderUserList(users) {
    userList.innerHTML = "";

    if (users.length === 0) {
      userList.innerHTML = "<p>No users found.</p>";
      return;
    }

    const table = document.createElement("table");
    table.classList.add("user-table");

    const thead = document.createElement("thead");
    const headerRow = document.createElement("tr");
    headerRow.innerHTML = `
      <th>Email</th>
      <th>Role</th>
      <th>Actions</th>
    `;
    thead.appendChild(headerRow);
    table.appendChild(thead);

    const tbody = document.createElement("tbody");

    users.forEach((user) => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td>${user.email}</td>
        <td>${user.role}</td>
        <td>
          <button class="btn edit-btn" data-id="${user.id}">Edit</button>
          <button class="btn delete-btn" data-id="${user.id}">Delete</button>
        </td>
      `;
      tbody.appendChild(row);
    });

    table.appendChild(tbody);
    userList.appendChild(table);

    // Attach event listeners for edit and delete buttons
    document
      .querySelectorAll(".edit-btn")
      .forEach((button) => button.addEventListener("click", handleEditUser));
    document
      .querySelectorAll(".delete-btn")
      .forEach((button) => button.addEventListener("click", handleDeleteUser));
  }

  // Handle Create User
  createUserForm.addEventListener("submit", async function (event) {
    event.preventDefault();
    const formData = new FormData(createUserForm);
    const data = {
      email: formData.get("email"),
      password: formData.get("password"),
      role: formData.get("role"),
    };

    try {
      const response = await fetch("http://localhost:4000/v1/user", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authToken}`,
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        const errorText = await response.text();
        createUserMessage.textContent = `Error: ${errorText}`;
        createUserMessage.style.display = "block";
        return;
      }

      createUserMessage.textContent = "User created successfully!";
      createUserMessage.style.display = "block";
      createUserForm.reset();
      fetchUsers();
    } catch (error) {
      createUserMessage.textContent = "An error occurred. Please try again.";
      createUserMessage.style.display = "block";
      console.error("Error creating user:", error);
    }
  });

  // Handle Edit User
  function handleEditUser(event) {
    const userId = event.target.getAttribute("data-id");

    // Fetch user data
    fetch(`http://localhost:4000/v1/user/${userId}`, {
      headers: {
        Authorization: `Bearer ${authToken}`,
      },
    })
      .then((response) => response.json())
      .then((user) => {
        document.getElementById("updateUserId").value = user.id;
        document.getElementById("updateEmail").value = user.email;
        document.getElementById("updateRole").value = user.role;
        updateUserSection.style.display = "block";
      })
      .catch((error) => console.error("Error fetching user:", error));
  }

  // Handle Update User
  updateUserForm.addEventListener("submit", async function (event) {
    event.preventDefault();
    const userId = document.getElementById("updateUserId").value;
    const formData = new FormData(updateUserForm);
    const data = {
      email: formData.get("email"),
      password: formData.get("password"),
      role: formData.get("role"),
    };

    try {
      const response = await fetch(`http://localhost:4000/v1/user/${userId}`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authToken}`,
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        const errorText = await response.text();
        updateUserMessage.textContent = `Error: ${errorText}`;
        updateUserMessage.style.display = "block";
        return;
      }

      updateUserMessage.textContent = "User updated successfully!";
      updateUserMessage.style.display = "block";
      updateUserForm.reset();
      updateUserSection.style.display = "none";
      fetchUsers();
    } catch (error) {
      updateUserMessage.textContent = "An error occurred. Please try again.";
      updateUserMessage.style.display = "block";
      console.error("Error updating user:", error);
    }
  });

  // Handle Cancel Update
  cancelUpdateButton.addEventListener("click", function () {
    updateUserSection.style.display = "none";
    updateUserForm.reset();
    updateUserMessage.style.display = "none";
  });

  // Handle Delete User
  function handleDeleteUser(event) {
    const userId = event.target.getAttribute("data-id");

    if (confirm("Are you sure you want to delete this user?")) {
      fetch(`http://localhost:4000/v1/user/${userId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("Failed to delete user");
          }
          fetchUsers();
        })
        .catch((error) => console.error("Error deleting user:", error));
    }
  }

  // Fetch users on page load
  fetchUsers();

  // Refresh user list
  refreshUserListButton.addEventListener("click", fetchUsers);
});
