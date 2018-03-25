$(document).ready(function() {

    $('.ui.radio.checkbox')
  	.checkbox();

    $('.ui.dropdown')
	.dropdown();

    $("#actType").change(function(){
		switch($(this).val())
		{
			case "2":	//email
				$("#weibo").hide();
				$("#email").show();
				$("#url").hide();
				$("#mobile_push").hide();
				break;
			case "3":
				$("#weibo").hide();
				$("#email").hide();
				$("#url").show();
				$("#mobile_push").hide();
				break;
			case "4":	//weibo
				$("#email").hide();
				$("#weibo").show();
				$("#url").hide();
				$("#mobile_push").hide();
				break;
			case "6":	//Mobile Push
				$("#email").hide();
				$("#weibo").hide();
				$("#url").hide();
				$("#mobile_push").show();
				break;
			default:
				break;
		}
	});

    $('.chooseUrl').change(function(e){
		e.preventDefault();
		 switch ($("input[name='chooseUrl']:checked").val())
		 {
			 case "1":
				$("span.chooseSwitch").show();
				$("span.writeUrl").hide();
				//$("#api_key").val(api_key);
				break;
			case "2":
				$("span.chooseSwitch").hide();
				$("span.writeUrl").show();
				$("#postUrl_2").next().html('');
				if($("#act_apt_key").val() == api_key){
					$("#postUrl_2").val("");
				}
			default:
				break;
		}

	});

});
