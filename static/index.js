
$(document).ready(() => {

    $("#create-form").submit((e) => {
      e.preventDefault();

      let formData = {}
      $(e.target).serializeArray().forEach((el) => {
        formData[el.name] = el.value;
      })

      var flags = [];
      $("input:checkbox[name=flag]:checked").each(function() { 
        flags.push($(this).val()); 
      });

      m.request({
        method: "POST",
        url: "/cmd/zt100.new-demo",
        body: {
            ID: window.OID,
            Name: formData["name"],
            Domain: formData["domain"],
            Color: formData["color"],
            Vertical: formData["vertical"],
            Features: flags,
        },
      }).then(() => location.reload());

      
      return false;
    })

    $("#field-domain").on("input", (e) => {
      $("#field-name").val(e.target.value.split(".")[0]);
      $("#field-logo").attr("src", "https://logo.clearbit.com/"+e.target.value);
      $("#field-logo").css("opacity", "1");
    })
  })