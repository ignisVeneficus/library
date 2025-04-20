function displaySeriesList(data){
    var page = $("#SeriesList");
    displayPagination(page,data.pagination, displaySeriesList);
 
	// count
	var countNode = $(".nrOfSeries",page);
	countNode.html(data.pagination.qty);

    var tbody = $(".tableList tbody", page);
	//tbody.empty();
	var newBody = $("<tbody></tbody>");
	$.each(data.result, function(key, value) {
		row = $("<tr></tr>").attr("id","id_"+value.id).data("id",value.id)
    
		if((key%4)>1){
			row.addClass("zebra");
		}
    
		col1 = $("<td class=\"number\">"+value.id+"</td>").appendTo(row);
		col2 = $("<td>"+value.name+"</td>").appendTo(row);
		
		col3 = $("<td>"+value.url+"</td>").appendTo(row);
		col4 = $("<td class=\"number\">"+value.books+"</td>").appendTo(row);
		
		col5=$("<td class=\"buttonCell\"></td>").appendTo(row);
		if(value.url!=""){
            btn1=$("<button class=\"button smallButton fa-solid fa-up-right-from-square\"></button>").appendTo(col5).on("click",function(event){
                windows.open(value.url,'_blank');
            });
        }
		btn2=$("<button class=\"button smallButton fa-solid fa-pencil\"></button>").appendTo(col5).on("click",function(event){
			editTrack(value.id);
		});
		btn3=$("<button class=\"button smallButton fa-solid fa-book\"></button>").appendTo(col5).on("click",function(event){
			loadBookListBySeriesId(value.id);
            changePage("BookList");
   		});
		
		row.appendTo(newBody);
		/*
		row.hover(function(event){
			boldMapLine(value["id"],tracklines);
		},function(event){
			normalMapLine(value["id"],tracklines);
		});
		*/
	});
	tbody.replaceWith(newBody);


}