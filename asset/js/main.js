function OnDeleteRecord(value) {
	myRet = confirm("削除？");
	if ( myRet == true ){
		data = {"deleteIdx": value }
		PostExec("/", data)
	}
 }

 function PostExec(action, data){
	 // フォームの生成
	var form = document.createElement("form");
	form.setAttribute("action", action);
	form.setAttribute("method", "post");
	form.style.display = "none";
	document.body.appendChild(form);
	// パラメタの設定
	if (data !== undefined) {
		for (var paramName in data) {
			var input = document.createElement('input');
			input.setAttribute('type', 'hidden');
			input.setAttribute('name', paramName);
			input.setAttribute('value', data[paramName]);
			form.appendChild(input);
		}
	}
	// submit
	form.submit();
 }
