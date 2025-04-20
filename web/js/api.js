urlRoot = "/api/"
bookUrl = urlRoot+"book"
authorUrl = urlRoot+"author"
seriesUrl = urlRoot+"series"


/********************
* API functions
*********************/

//generic
function loadPagedData(pageUrl,page,callbackFv){
	console.debug(pageUrl);
	$.ajax({
		url			: pageUrl+"&page="+page,
		dataType	: 'json',
		data		: null,
		type		: 'GET',
		success		: function(jsonData){
			callbackFv(jsonData);
		},
		error		: function (xhr, ajaxOptions, thrownError) {
			alert(xhr.status);
			alert(thrownError);
		}
	});

}

// for books
function loadAllBookList(){
	var pageUrl = bookUrl+"?";
	loadPagedData(pageUrl,1,displayBookList);
}
// for all autors
function loadAllAuthorList(){
	var pageUrl = authorUrl+"?";
	loadPagedData(pageUrl,1,displayAuthorList);
}
// for all series
function loadAllSeriesList(){
	var pageUrl = seriesUrl+"?";
	loadPagedData(pageUrl,1,displaySeriesList);
}
// for all tags
function loadAllTagList(){
	var pageUrl = urlRoot+"tag?";
	loadPagedData(pageUrl,1,displayTagList);
}
function loadABook(id,callbackFv){
	$.ajax({
		url			: bookUrl+"/"+id,
		dataType	: 'json',
		data		: null,
		type		: 'GET',
		success		: function(jsonData){
			callbackFv(jsonData);//displayABook(jsonData);
		},
		error		: function (xhr, ajaxOptions, thrownError) {
			alert(xhr.status);
			alert(thrownError);
		}
	});
}
function loadBookListByAuthorId(authorId){
	var pageUrl = bookUrl+"?ai="+authorId;
	loadPagedData(pageUrl,1,displayBookList);
}
function loadBookListBySeriesId(seriesId){
	var pageUrl = bookUrl+"?si="+seriesId;
	loadPagedData(pageUrl,1,displayBookList);
}
function loadBookListByTagId(tagId){
	var pageUrl = bookUrl+"?ti="+tagId;
	loadPagedData(pageUrl,1,displayBookList);
}
function scrapeABook(url){
	$.ajax({
		url			: urlRoot+"scraper/?url="+encodeURIComponent(url),
		dataType	: 'json',
		data		: null,
		type		: 'GET',
		success		: function(jsonData){
			doPopupScraperBook(jsonData)
		},
		error		: function (xhr, ajaxOptions, thrownError) {
			alert(xhr.status);
			alert(thrownError);
		}
	});
}
function pushBook(book){
	$.ajax({
		url			: bookUrl,
		jsonp		: false,
		dataType	: 'json',
		data		: JSON.stringify(book),
		type		: 'POST',
		contentType	: "application/json",
		crossDomain	: true,
		success		: function(jsonData){
			changePage("BookList");
			reloadPage("#BookList");
		},
		error		: function (xhr, ajaxOptions, thrownError) {
			alert(xhr.status);
			alert(thrownError);
		}
	});
}

