document.addEventListener("DOMContentLoaded", function () {
  const viewUserForm = document.getElementById("viewUserForm");
  const viewUserMessage = document.getElementById("viewUserMessage");
  const userDetail = document.getElementById("userDetail");

  let authToken = localStorage.getItem("authToken");

  // Handle View User
  viewUserForm.addEventListener("submit", async function (event) {
    event.preventDefault();
    const userId = document.getElementById("viewUserId").value;

    try {
      const response = await fetch(`http://localhost:4000/v1/user/${userId}`, {
        mode: "cors",
      });

      if (!response.ok) {
        throw new Error(`User not found with ID: ${userId}`);
      }

      const user = await response.json();
      renderUserDetail(user);
    } catch (error) {
      viewUserMessage.textContent = error.message;
      viewUserMessage.style.display = "block";
      userDetail.innerHTML = "";
      console.error("Error fetching user:", error);
    }
  });

  // Render User Detail
  function renderUserDetail(user) {
    viewUserMessage.style.display = "none";
    userDetail.innerHTML = `
      <h2>User Details</h2>
      <p><strong>Name:</strong> ${user.name}</p>
      <p><strong>Email:</strong> ${user.email}</p>
      <p><strong>Role:</strong> ${user.role}</p>
      <p><strong>Created At:</strong> ${new Date(
        user.created_at
      ).toLocaleDateString()}</p>
    `;
  }
});
