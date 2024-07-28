document.addEventListener("DOMContentLoaded", function () {
  const payrollList = document.getElementById("payrollList");
  let authToken = localStorage.getItem("authToken");

  // Fetch Payrolls
  async function fetchPayrolls() {
    payrollList.innerHTML = "Loading...";
    try {
      const response = await fetch("http://localhost:4000/v1/payroll", {
        headers: {
          "Content-Type": "application/json",
        },
        mode: "cors",
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      renderPayrollList(data);
    } catch (error) {
      payrollList.innerHTML = "Error loading payrolls";
      console.error("Error fetching payrolls:", error);
    }
  }

  // Render Payroll List
  function renderPayrollList(payrolls) {
    payrollList.innerHTML = "";

    if (payrolls.length === 0) {
      payrollList.innerHTML = "<p>No payrolls found.</p>";
      return;
    }

    const table = document.createElement("table");
    table.classList.add("payroll-table");

    const thead = document.createElement("thead");
    const headerRow = document.createElement("tr");
    headerRow.innerHTML = `
            <th>Employee ID</th>
            <th>Amount</th>
            <th>Date</th>
        `;
    thead.appendChild(headerRow);
    table.appendChild(thead);

    const tbody = document.createElement("tbody");

    payrolls.forEach((payroll) => {
      const row = document.createElement("tr");
      row.innerHTML = `
                <td>${payroll.employee_id}</td>
                <td>${payroll.amount}</td>
                <td>${new Date(payroll.date).toLocaleDateString()}</td>
            `;
      tbody.appendChild(row);
    });

    table.appendChild(tbody);
    payrollList.appendChild(table);

    // Attach event listeners for edit and delete buttons
    document
      .querySelectorAll(".edit-btn")
      .forEach((button) => button.addEventListener("click", handleEditPayroll));
    document
      .querySelectorAll(".delete-btn")
      .forEach((button) =>
        button.addEventListener("click", handleDeletePayroll)
      );
  }

  // Handle Edit Payroll
  function handleEditPayroll(event) {
    const payrollId = event.target.getAttribute("data-id");
    // Redirect to edit payroll page (not implemented in this example)
    window.location.href = `edit_payroll.html?id=${payrollId}`;
  }

  // Handle Delete Payroll
  function handleDeletePayroll(event) {
    const payrollId = event.target.getAttribute("data-id");

    if (confirm("Are you sure you want to delete this payroll?")) {
      fetch(`http://localhost:4000/v1/payroll/${payrollId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${authToken}`,
          "Content-Type": "application/json",
        },
      })
        .then((response) => {
          if (response.ok) {
            alert("Payroll deleted successfully!");
            fetchPayrolls(); // Refresh the list
          } else {
            alert("Failed to delete payroll.");
          }
        })
        .catch((error) => {
          console.error("Error deleting payroll:", error);
        });
    }
  }

  // Initial Fetch of Payrolls
  fetchPayrolls();
});
