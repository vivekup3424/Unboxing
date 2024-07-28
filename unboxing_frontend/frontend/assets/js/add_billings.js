document.addEventListener("DOMContentLoaded", function () {
  const addBillingForm = document.getElementById("addBillingForm");
  const addBillingMessage = document.getElementById("addBillingMessage");
  let authToken = localStorage.getItem("authToken");

  addBillingForm.addEventListener("submit", async function (event) {
    event.preventDefault();

    // Collect form data
    const formData = new FormData(addBillingForm);
    const data = {
      customer_id: parseInt(formData.get("customer_id")),
      amount: parseFloat(formData.get("amount")),
      date: formData.get("date"),
    };

    try {
      const response = await fetch("http://localhost:4000/v1/billing", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
        mode: "cors",
      });

      if (!response.ok) {
        addBillingMessage.textContent = `Error: ${
          responseData.message || "Unable to add billing"
        }`;
        addBillingMessage.style.display = "block";
        return;
      }

      const responseData = await response.json();
      addBillingMessage.textContent = "Billing added successfully!";
      addBillingMessage.style.display = "block";
      addBillingForm.reset();
    } catch (error) {
      addBillingMessage.textContent = "Billing added successfully!";
      addBillingMessage.style.display = "block";
      console.error("Error adding billing:", error);
    }
  });
});
