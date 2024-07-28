document.addEventListener("DOMContentLoaded", function () {
  const customerList = document.getElementById("customerList");
  const refreshCustomerList = document.getElementById("refreshCustomerList");

  let authToken = localStorage.getItem("authToken");

  // Fetch All Customers
  async function fetchCustomers() {
    try {
      const response = await fetch("http://localhost:4000/v1/customer", {
        method: "GET",
      });

      if (!response.ok) {
        customerList.innerHTML = "<p>Error loading customers.</p>";
        return;
      }

      const customers = await response.json();
      renderCustomerList(customers);
    } catch (error) {
      customerList.innerHTML = "<p>An error occurred. Please try again.</p>";
      console.error("Error fetching customers:", error);
    }
  }

  // Render Customer List
  function renderCustomerList(customers) {
    if (customers.length === 0) {
      customerList.innerHTML = "<p>No customers found.</p>";
      return;
    }

    let html = '<table class="customer-table">';
    html += "<tr><th>ID</th><th>Name</th><th>Email</th><th>Phone</th></tr>";
    customers.forEach((customer) => {
      html += `<tr>
        <td>${customer.id}</td>
        <td>${customer.name}</td>
        <td>${customer.email}</td>
        <td>${customer.phone}</td>
      </tr>`;
    });
    html += "</table>";

    customerList.innerHTML = html;
  }

  // Initial Fetch
  fetchCustomers();

  // Refresh List
  refreshCustomerList.addEventListener("click", fetchCustomers);
});
