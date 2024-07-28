document.addEventListener("DOMContentLoaded", function () {
  const deleteCustomerForm = document.getElementById("deleteCustomerForm");
  const deleteCustomerMessage = document.getElementById(
    "deleteCustomerMessage"
  );

  let authToken = localStorage.getItem("authToken");

  // Handle Delete Customer
  deleteCustomerForm.addEventListener("submit", async function (event) {
    event.preventDefault();
    const customerId = document.getElementById("deleteCustomerId").value;

    try {
      const response = await fetch(
        `http://localhost:4000/v1/customer/${customerId}`,
        {
          method: "DELETE",
        }
      );

      if (!response.ok) {
        const errorText = await response.text();
        deleteCustomerMessage.textContent = `Error: ${errorText}`;
        deleteCustomerMessage.style.display = "block";
        return;
      }

      deleteCustomerMessage.textContent = "Customer deleted successfully!";
      deleteCustomerMessage.style.display = "block";
      deleteCustomerForm.reset();
    } catch (error) {
      deleteCustomerMessage.textContent =
        "An error occurred. Please try again.";
      deleteCustomerMessage.style.display = "block";
      console.error("Error deleting customer:", error);
    }
  });
});
