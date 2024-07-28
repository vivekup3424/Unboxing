document.addEventListener("DOMContentLoaded", function () {
  const userList = document.getElementById("userList");
  const refreshUserListButton = document.getElementById("refreshUserList");

  let authToken = localStorage.getItem("authToken");

  // Fetch Users
  async function fetchUsers() {
    userList.innerHTML = "Loading...";
    try {
      const response = await fetch("http://localhost:4000/v1/user", {
        method: "GET",
        mode: "cors",
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      console.log("Response:", response);
      const data = await response.json();
      console.log(data);
      renderUserList(data);
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
      <th>Name</th>
      <th>Email</th>
      <th>Role</th>
      <th>Created At</th>
    `;
    thead.appendChild(headerRow);
    table.appendChild(thead);

    const tbody = document.createElement("tbody");
    //iterating over all the users
    users.forEach((user) => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td>${user.name}</td>
        <td>${user.email}</td>
        <td>${user.role}</td>
        <td>${new Date(user.created_at).toLocaleDateString()}</td>
      `;
      tbody.appendChild(row);
    });

    table.appendChild(tbody);
    userList.appendChild(table);
  }

  // Fetch users on page load
  fetchUsers();

  // Refresh user list
  refreshUserListButton.addEventListener("click", fetchUsers);
});
