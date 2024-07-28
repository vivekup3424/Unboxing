document.addEventListener("DOMContentLoaded", function () {
  const createUserForm = document.getElementById("createUserForm");
  const createUserMessage = document.getElementById("createUserMessage");

  let authToken = localStorage.getItem("authToken");

  // Handle Create User
  createUserForm.addEventListener("submit", async function (event) {
    event.preventDefault();
    const formData = new FormData(createUserForm);
    const data = {
      name: formData.get("name"),
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
    } catch (error) {
      createUserMessage.textContent = "An error occurred. Please try again.";
      createUserMessage.style.display = "block";
      console.error("Error creating user:", error);
    }
  });
});
