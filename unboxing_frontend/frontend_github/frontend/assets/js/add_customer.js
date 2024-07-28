document.addEventListener("DOMContentLoaded", function () {
  const addCustomerForm = document.getElementById("addCustomerForm");
  const addCustomerMessage = document.getElementById("addCustomerMessage");

  let authToken = localStorage.getItem("authToken");

  // Handle Add Customer
  addCustomerForm.addEventListener("submit", async function (event) {
    event.preventDefault();
    const formData = new FormData(addCustomerForm);
    const data = {
      name: formData.get("name"),
      email: formData.get("email"),
      phone: formData.get("phone"),
    };

    try {
      const response = await fetch("http://localhost:4000/v1/customer", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authToken}`,
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        const errorText = await response.text();
        addCustomerMessage.textContent = `Error: ${errorText}`;
        addCustomerMessage.style.display = "block";
        return;
      }

      addCustomerMessage.textContent = "Customer added successfully!";
      addCustomerMessage.style.display = "block";
      addCustomerForm.reset();
    } catch (error) {
      addCustomerMessage.textContent = "An error occurred. Please try again.";
      addCustomerMessage.style.display = "block";
      console.error("Error adding customer:", error);
    }
  });
});
