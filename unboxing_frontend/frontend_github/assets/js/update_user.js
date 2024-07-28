document.addEventListener("DOMContentLoaded", function () {
  const updateUserForm = document.getElementById("updateUserForm");
  const updateUserMessage = document.getElementById("updateUserMessage");
  const cancelUpdateButton = document.getElementById("cancelUpdate");

  let authToken = localStorage.getItem("authToken");

  // Handle Update User
  updateUserForm.addEventListener("submit", async function (event) {
    event.preventDefault();
    const formData = new FormData(updateUserForm);
    const data = {
      email: formData.get("email"),
      password: formData.get("password"),
      role: formData.get("role"),
    };

    const userId = formData.get("userId");

    try {
      const response = await fetch(`http://localhost:4000/v1/user/${userId}`, {
        method: "PUT",
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
    } catch (error) {
      updateUserMessage.textContent = "An error occurred. Please try again.";
      updateUserMessage.style.display = "block";
      console.error("Error updating user:", error);
    }
  });

  // Handle Cancel Update
  cancelUpdateButton.addEventListener("click", function () {
    updateUserForm.reset();
    updateUserMessage.style.display = "none";
  });
});
