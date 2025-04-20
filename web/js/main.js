function downloadBook(){
	var file = $('#BookOverlay').data("bookFile");
	window.open('books/'+file, '_blank');
}
function searchMolyBook(){
	var title = $("#BookEdit_title").val();
	var title = title.replaceAll(' ', '+');
	window.open('https://moly.hu/kereses?utf8=%E2%9C%93&query='+title, '_blank');	
}
function openBookLink(){
	var url = $("#BookOverlay").data("url");
	if((url!=null)&&(url!="")){
		window.open(url, '_blank');	
	}
}

function filterAuthor(){
	var url = authorUrl+"?";
	var data = $("#AuthorList_search").val();
	if(data.trim().length !== 0){
		url+="q="+encodeURI(data);
	}
	loadPagedData(url,1,displayAuthorList);
}

function filterSeries(){
	var url = seriesUrl+"?";
	var data = $("#SeriesList_search").val();
	if(data.trim().length !== 0){
		url+="q="+encodeURI(data);
	}
	loadPagedData(url,1,displaySeriesList);
}

function filterTag(){
	var url = tagUrl+"?";
	var data = $("#TagList_search").val();
	if(data.trim().length !== 0){
		url+="q="+encodeURI(data);
	}
	loadPagedData(url,1,displayTagList);
}

function filterBook(){
	var url = bookUrl+"?";
	var data = $("#BookList_search").val();
	if(data.trim().length !== 0){
		url+="q="+encodeURI(data);
	}
	loadPagedData(url,1,displayBookList);
}

function addData(btn,nr,baseUrl,callbackFv){
	btn.data({page:nr, baseUrl: baseUrl,callback:callbackFv});
}

function fillFilters(filterNode,filters,loadFv){
	filterNode.empty()
	if((filters!=null) && (filters.length >0)){
		$.each(filters, function(key, filter) {
			var node = $("<div class=\"filter\"></div>");
			node.appendTo(filterNode);
			var text = $("<div class=\"text\">Filtered by " + filter.type +": " +filter.value +"</div>");
			text.appendTo(node);
			var button = $("<button class=\"button smallButton fa-solid fa-xmark\"></button>").appendTo(node);
			button.on("click", function(e){
				loadFv();
			});
		});
	}

}

function reloadPage(page){
	var node = $(page);
	var div = $(".pagination",node);
	var data = div.data("data");
	var callbackFv = div.data("callback");
	loadPagedData(data.base,data.selectedPage,callbackFv);
}
function displayPagination(node,data,callbackFv){
	var div = $(".pagination",node);
	div.data("data",data);
	div.data("callback",callbackFv);

	var nrOfPages = data.pages;
	var selected = data.selectedPage;
	var baseUrl = data.base;
	var buttons = $(".pagesButtons",node);
	var pageNode = $("input",buttons);
	pageNode.val(selected);

	var nrNode = $(".maxPage",node);
	nrNode.html("of "+nrOfPages);
	
	var prev = $(".prevButton",buttons);
	var first = $(".firstButton",buttons);
	var next = $(".nextButton",buttons);
	var last = $(".lastButton",buttons);
	var jump = $(".jump",buttons);

	var isFirst = (selected ==1);
	var isLast = (selected==nrOfPages);
	prev.prop('disabled', isFirst);
	first.prop('disabled', isFirst);
	next.prop('disabled', isLast);
	last.prop('disabled', isLast);

	var reload=$(".reload",buttons);


	addData(prev,selected-1,baseUrl,callbackFv);
	addData(first,1,baseUrl,callbackFv);
	addData(next,selected+1,baseUrl,callbackFv);
	addData(last,nrOfPages,baseUrl,callbackFv);

	addData(reload,selected,baseUrl,callbackFv);

	pageNode.data({max:nrOfPages, baseUrl: baseUrl,callback:callbackFv,button:jump}).on( "focusout",function(e){
		var val = $(this).val();
		val = parseInt(+val);
		var max = $(this).data("max");
		if((val == NaN) || (val<1)){
			val = 1;
		}
		if( val > max){
			val = max;
		}
		var button= $(this).data("button");
		button.data("page",val);
	});
	addData(jump,selected,baseUrl,callbackFv);

}

function gotoPage(node){
	let dest = node.data("page");
	let baseUrl = node.data("baseUrl");
	let callbackFv = node.data("callback");
	loadPagedData(baseUrl,dest,callbackFv);
}

function searchEnter(e){
	if (e.key==="Enter") {
		let element = document.activeElement;
		let node = $(element);
		if(node.attr("name")=="search"){
			let id = node.attr('id');
			$("#"+id+"_btn").click();
			e.preventDefault();
			e.stopPropagation();
		}
	}

}

function init(){

	$(".menuItem .item").on('click',function(event){
		changePage($(this).attr("panel"));
	});
	$('#BookOverlay_download').on("click",function(event){
		downloadBook();
	});
	$('#BookEdit_search').on("click",function(event){
		searchMolyBook();
	});
	$('#BookEdit_srape').on("click",function(event){
		scrapeBook();
	});
	$("#BookEdit_cancel").on('click',function(event){
		changePage("BookList");
	});
	$("#BookEdit_save").on('click',function(event){
		var book = readBackBook();
		pushBook(book);
	});
	$("#BookScraper_title").on('click',function(event){
		scraperMoveStatic("#scraperTitleFrom","#scraperTitleTo");
	})
	$("#BookScraper_blurb").on('click',function(event){
		scraperMoveStatic("#scraperBlurbFrom","#scraperBlurbTo");
	})
	$("#BookScraper_replaceAuthor").on('click',function(event){
		scraperMoveList("#scraperAuthorFrom","#scraperAuthorTo");
	})
	$("#BookScraper_replaceSeries").on('click',function(event){
		scraperMoveList("#scraperSeriesFrom","#scraperSeriesTo");
	})
	$("#BookScraper_replaceTag").on('click',function(event){
		scraperMoveList("#scraperTagFrom","#scraperTagTo");
	})

	$("#BookScraper_addAuthor").on('click',function(event){
		scraperAddList("#scraperAuthorFrom",addScraperAuthor);
	})
	$("#BookScraper_addSeries").on('click',function(event){
		scraperAddList("#scraperSeriesFrom",addScraperSeries);
	})
	$("#BookScraper_addTag").on('click',function(event){
		scraperAddList("#scraperTagsFrom",addScraperTag);
	})

	$('#BookOverlay_edit').on("click",function(event){
		editBook(event);		
	});
	$('#BookOverlay_link').on('click',function(event){
		openBookLink();
	});
	$("#BookOverlay_cancel").on('click',function(event){
		$('#BookOverlay').click();
	});

	$("#BookScraper_cancel").on('click',function(event){
		$('#BookScraperOverlay').click();
	});
	$("#BookScraper_ok").on('click',function(event){
		writeBackScraper();
	});

	$(".pagination button").on("click", function(event){
		gotoPage($(this));
	});

	$('#BookList_search_btn').on("click",function(event){
		filterBook();
	});

	$('#AuthorList_search_btn').on("click",function(event){
		filterAuthor();
	});

	$('#SeriesList_search_btn').on("click",function(event){
		filterSeries();
	});


	$('button[title]').addTooltip();



	loadAllBookList();
	loadAllAuthorList();
	loadAllSeriesList();

	$(document).on("keypress",searchEnter);
}

function changePage(newId){
	$("#"+newId).removeClass("hidden");
	$(".page").not("#"+newId).addClass("hidden");
	$("#menuitem-"+newId).addClass("selected");
	$(".menuItem .item").not("#menuitem-"+newId).removeClass("selected");
}
