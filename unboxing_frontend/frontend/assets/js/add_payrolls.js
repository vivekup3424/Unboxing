document.addEventListener("DOMContentLoaded", function () {
  const addPayrollForm = document.getElementById("addPayrollForm");
  const addPayrollMessage = document.getElementById("addPayrollMessage");
  let authToken = localStorage.getItem("authToken");

  // Handle Add Payroll Form Submission
  addPayrollForm.addEventListener("submit", async function (event) {
    event.preventDefault();

    // Collect form data
    const formData = new FormData(addPayrollForm);
    const data = {
      employee_id: parseInt(formData.get("employee_id")),
      amount: parseFloat(formData.get("amount")),
      date: formData.get("date"),
    };

    try {
      const response = await fetch("http://localhost:4000/v1/payroll", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
        mode: "cors",
      });

      if (!response.ok) {
        const errorText = await response.text();
        addPayrollMessage.textContent = `Error: ${errorText}`;
        addPayrollMessage.style.display = "block";
        return;
      }

      const result = await response.json();
      addPayrollMessage.textContent = "Payroll added successfully!";
      addPayrollMessage.style.display = "block";
      addPayrollForm.reset();
    } catch (error) {
      addPayrollMessage.textContent = "An error occurred. Please try again.";
      addPayrollMessage.style.display = "block";
      console.error("Error adding payroll:", error);
    }
  });
});
