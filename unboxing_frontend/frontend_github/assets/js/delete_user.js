document.addEventListener("DOMContentLoaded", function () {
  const deleteUserForm = document.getElementById("deleteUserForm");
  const deleteUserMessage = document.getElementById("deleteUserMessage");

  let authToken = localStorage.getItem("authToken");

  // Handle Delete User
  deleteUserForm.addEventListener("submit", async function (event) {
    event.preventDefault();
    const userId = document.getElementById("deleteUserId").value;

    try {
      const response = await fetch(`http://localhost:4000/v1/user/${userId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });

      if (!response.ok) {
        const errorText = await response.text();
        deleteUserMessage.textContent = `Error: ${errorText}`;
        deleteUserMessage.style.display = "block";
        return;
      }

      deleteUserMessage.textContent = "User deleted successfully!";
      deleteUserMessage.style.display = "block";
      deleteUserForm.reset();
    } catch (error) {
      deleteUserMessage.textContent = "An error occurred. Please try again.";
      deleteUserMessage.style.display = "block";
      console.error("Error deleting user:", error);
    }
  });
});
