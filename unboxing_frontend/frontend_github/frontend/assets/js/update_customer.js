document.addEventListener("DOMContentLoaded", function () {
  const updateCustomerForm = document.getElementById("updateCustomerForm");
  const updateCustomerMessage = document.getElementById(
    "updateCustomerMessage"
  );
  const updateCustomerId = document.getElementById("updateCustomerId");
  const cancelUpdateButton = document.getElementById("cancelUpdate");

  let authToken = localStorage.getItem("authToken");

  // Handle Update Customer
  updateCustomerForm.addEventListener("submit", async function (event) {
    event.preventDefault();
    const formData = new FormData(updateCustomerForm);
    const data = {
      email: formData.get("email"),
      phone: formData.get("phone"),
    };

    const customerId = updateCustomerId.value;

    try {
      const response = await fetch(
        `http://localhost:4000/v1/customer/${customerId}`,
        {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${authToken}`,
          },
          body: JSON.stringify(data),
        }
      );

      if (!response.ok) {
        const errorText = await response.text();
        updateCustomerMessage.textContent = `Error: ${errorText}`;
        updateCustomerMessage.style.display = "block";
        return;
      }

      updateCustomerMessage.textContent = "Customer updated successfully!";
      updateCustomerMessage.style.display = "block";
      updateCustomerForm.reset();
    } catch (error) {
      updateCustomerMessage.textContent =
        "An error occurred. Please try again.";
      updateCustomerMessage.style.display = "block";
      console.error("Error updating customer:", error);
    }
  });

  // Handle Cancel Update
  cancelUpdateButton.addEventListener("click", function () {
    updateCustomerForm.reset();
    updateCustomerMessage.style.display = "none";
  });
});
