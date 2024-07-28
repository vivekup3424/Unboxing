document.addEventListener("DOMContentLoaded", function () {
  const billingsTableBody = document.querySelector("#billingsTable tbody");
  const billingMessage = document.getElementById("billingMessage");
  let authToken = localStorage.getItem("authToken");

  async function fetchBillings() {
    try {
      const response = await fetch("http://localhost:4000/v1/billing", {
        method: "GET",
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
        mode: "cors",
      });

      if (!response.ok) {
        billingMessage.textContent = `Error: ${
          data.message || "Unable to load billings"
        }`;
        billingMessage.style.display = "block";
        return;
      }

      const data = await response.json();
      console.log(data);
      renderBillings(data);
    } catch (error) {
      billingMessage.textContent = "An error occurred while fetching billings.";
      billingMessage.style.display = "block";
      console.error("Error fetching billings:", error);
    }
  }

  function renderBillings(billings) {
    billingsTableBody.innerHTML = "";
    billings.forEach((billing) => {
      const row = document.createElement("tr");
      row.innerHTML = `
                <td>${billing.id}</td>
                <td>${billing.customer_id}</td>
                <td>${billing.amount.toFixed(2)}</td>
                <td>${new Date(billing.date).toLocaleDateString()}</td>
            `;
      billingsTableBody.appendChild(row);
    });
  }

  fetchBillings();
});
