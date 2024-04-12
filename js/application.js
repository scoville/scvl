document.addEventListener("DOMContentLoaded", function () {
  // data-confirm
  document.querySelectorAll("[data-confirm]").forEach(function (element) {
    element.addEventListener("click", function (e) {
      if (!confirm(e.target.getAttribute("data-confirm"))) {
        e.preventDefault();
        return false;
      }
    });
  });

  if (document.getElementById("shortenUrl")) {
    document.getElementById("shortenUrl").value = location.origin + "/" + document.getElementById("shortenUrl").value;
  }
  const $copy = document.getElementById("copy");
  if ($copy) {
    $copy.addEventListener("click", function () {
      const copyTarget = $copy.getAttribute("data-copy-target");
      document.querySelector(copyTarget).select();
      document.execCommand("copy");
      $copy.querySelector(".hint").style.visibility = "visible";
      $copy.querySelector(".hint").textContent = "コピー完了";
    });
  }

  // UTM
  if (document.getElementById("url")) {
    var params = getParams(document.getElementById("url").value);
    if (params["utm_source"]) {
      document.querySelector('input[name="utm_source"]').value = params["utm_source"];
    }
    if (params["utm_medium"]) {
      document.querySelector('select[name="utm_medium"]').value = params["utm_medium"];
    }
    if (params["utm_campaign"]) {
      document.querySelector('input[name="utm_campaign"]').value = params["utm_campaign"];
    }
  }
  var utmCheckbox = document.getElementById("utm");
  if (utmCheckbox) {
    if (utmCheckbox.checked) {
      document.getElementById("utm-fields").style.display = "block";
    }
    utmCheckbox.addEventListener("change", function () {
      if (this.checked) {
        document.getElementById("utm-fields").style.display = "block";
      } else {
        document.getElementById("utm-fields").style.display = "none";
      }
    });
  }
  document.querySelectorAll('input[name="url"], input[name="utm_source"], select[name="utm_medium"], input[name="utm_campaign"]').forEach(function (element) {
    element.addEventListener("change", buildURL);
  });

  // OGP
  var ogpCheckbox = document.getElementById("ogp");
  if (ogpCheckbox) {
    if (ogpCheckbox.checked) {
      document.getElementById("ogp-fields").style.display = "block";
    }
    ogpCheckbox.addEventListener("change", function () {
      if (this.checked) {
        document.getElementById("ogp-fields").style.display = "block";
        if (document.querySelector('input.ogp[name="title"]').value === "") {
          fetchOGP();
        }
      } else {
        document.getElementById("ogp-fields").style.display = "none";
      }
    });
  }

  // Email
  var emailCheckbox = document.getElementById("email");
  if (emailCheckbox) {
    if (emailCheckbox.checked) {
      document.getElementById("email-fields").style.display = "block";
    }
    emailCheckbox.addEventListener("change", function () {
      if (this.checked) {
        document.getElementById("email-fields").style.display = "block";
      } else {
        document.getElementById("email-fields").style.display = "none";
      }
    });
  }

  function getParams(url) {
    var params = {};
    if (url.split("?").length <= 1) {
      return params;
    }
    var q = url.split("?")[1];
    q.split("&").forEach(function (part) {
      var param = part.split("=");
      if (param.length > 1) {
        params[param[0]] = param[1];
      }
    });
    return params;
  }

  function buildURL() {
    var urlField = document.getElementById("url");
    var url = urlField.value;
    if (url === "") {
      return;
    }
    var params = getParams(url);
    var source = document.querySelector('input[name="utm_source"]').value;
    if (source) {
      params["utm_source"] = source;
    }
    var medium = document.querySelector('select[name="utm_medium"]').value;
    if (medium === "-") {
      delete params["utm_medium"];
    } else if (medium) {
      params["utm_medium"] = medium;
    }
    var campaign = document.querySelector('input[name="utm_campaign"]').value;
    if (campaign) {
      params["utm_campaign"] = campaign;
    }
    url = url.split("?")[0];
    var values = [];
    Object.keys(params).forEach(function (key) {
      var val = params[key];
      if (val !== "") {
        values.push(key + "=" + val);
      }
    });
    if (values.length > 0) {
      url += "?" + values.join("&");
    }
    urlField.value = url;
  }

  function fetchOGP() {
    var url = document.getElementById("url").value;
    if (url === "") {
      return;
    }
    document.querySelectorAll("input.ogp").forEach(function (element) {
      element.value = "取得中...";
      element.disabled = true;
    });
    var request = new XMLHttpRequest();
    request.open("GET", "https://ogp.en-courage.com?url=" + url, true);
    request.responseType = 'json';
    request.onload = function () {
      var data = request.response;
      document.querySelector('input.ogp[name="title"]').value = data.title;
      document.querySelector('input.ogp[name="image"]').value = data.image;
      document.querySelector('input.ogp[name="description"]').value = data.description;
      document.querySelectorAll("input.ogp").forEach(function (element) {
        element.disabled = false;
      });
    };
    request.send();
  }
});