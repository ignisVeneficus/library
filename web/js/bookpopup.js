
function doPopupBook(book){
	var focused = $(':focus');
	var prevZ = maxZ;
	maxZ = maxZ+10;
	var thisZ = maxZ;
	var overlay = $('#BookOverlay');
    overlay.data('focus',focused).data('prevZ',prevZ).data('bookId',book.id).data('bookFile',book.file).data('url',book.url);
    var box =$('#BookOverlay .popup');

	if((book.url!=null)&&(book.url!="")){
		$("#BookOverlay_link").show();
	}
	else{
		$("#BookOverlay_link").hide();
	}

    $(".title",overlay).html(book.title)
    var authors = [];
    for( var i=0;i<book.authors.length;i++){
        authors.push(book.authors[i].name);
    }
    var authorName = authors.join("; ");
    var coverNode =  $(".cover",overlay);
    if(book.hasCover){
       coverNode.css({backgroundImage:"url(cover/"+book.file+"."+book.coverType+")",backgroundColor:""})
       coverNode.empty();
    }
    else{
        color = colors[book.id%colors.length];
        coverNode.css({backgroundImage:"",backgroundColor:color})
        var c = hexToRgb(color);
        createBlankCoverContent(coverNode,setContrast(c),book.title, authorName);
    }
	authorContainer = $(".authorContainer",overlay);
	authorContainer.empty();
	$.each(book.authors, function(key, author) {
		var line = $("<div></div>").addClass("lineContent").appendTo(authorContainer);
		var auth = $("<div class=\"author\">"+author.name+"</div>").appendTo(line);
		var buttons = $("<div class=\"buttons\"></div>").appendTo(line);
		var button = $("<button class=\"button smallButton fa-solid fa-book\" title=\"Show author's books\"></button>").appendTo(buttons).addTooltip();
		button.on("click", function(e){
			overlay.click();
			loadBookListByAuthorId(author.id);
		});
		if((author.url!=null)&&(author.url!="")){
			var button = $("<button class=\"button smallButton fa-solid fa-up-right-from-square\" title=\"Go to author's webpage\"></button>").appendTo(buttons).addTooltip();
			button.on("click", function(e){
				window.open(author.url, '_blank');
			});
		}

	});

	var seriesContainer = $(".seriesContainer",overlay);
	seriesContainer.empty();
	$.each(book.series, function(key, series) {
		var line = $("<div></div>").addClass("lineContent").appendTo(seriesContainer);
		var ser = $("<div class=\"series\"></div>").appendTo(line);
		var txt = series.name;
		if((series.seqno!=null)&&(series.seqno!=0)){
			txt +=" #"+series.seqno;
		}
		ser.html(txt);
		var buttons = $("<div class=\"buttons\"></div>").appendTo(line);
		var button = $("<button class=\"button smallButton fa-solid fa-book\" title=\"Show books in series\"></button>").appendTo(buttons).addTooltip();
		button.on("click", function(e){
			overlay.click();
			loadBookListBySeriesId(series.id);
		});
		if((series.url!=null)&&(series.url!="")){
			var button = $("<button class=\"button smallButton fa-solid fa-up-right-from-square\" title=\"Go to series's webpage\"></button>").appendTo(buttons).addTooltip();
			button.on("click", function(e){
				window.open(series.url, '_blank');
			});
		}

	});

	var tagsContainer = $(".tagContainer",overlay);
	tagsContainer.empty();
	$.each(book.tags, function(key, tag) {
		var tagN = $("<button class=\"tag\" title=\"Show books with this tag\"></button>").appendTo(tagsContainer).addTooltip();;
		tagN.html(tag.name);
		var color = tag.color;
		if((color==null) || (color=="")){
			color = colors[tag.id%colors.length];
		}
		var c = hexToRgb(color);
		var font = setContrast(c)
		tagN.css({color:font,backgroundColor:color});
		tagN.on("click", function(e){
			overlay.click();
			loadBookListByTagId(tag.id);
		});
	});

	var blurbContainer = $(".blurb",overlay);
	var txt = "";
	if(book.blurb!=null){
		txt = book.blurb.replaceAll("\n","<br>")
	}
	blurbContainer.html(txt);

	var filecontainer = $(".file",overlay);
	filecontainer.html(book.file)


	$(document).on("keypress.closeBook", function(e) {
		if (e.key === "Escape" || e.key==="Enter") {
			overlay.click();
		}
		if (e.key === "e") {
			editBook(e);
		}
		e.preventDefault();
		e.stopPropagation();
	});
    overlay.show();

	focused.blur();
	overlay.show().on("click",function(e){
		if (e.target != e.currentTarget) return;
		closePopupBook(e,'#BookOverlay')
	});
	//resizePopup(box,null);
	//preventScroll(overlay);
	//spreventScroll(box);
}

function closePopupBook(e,overlayId){
	var overlay = $(overlayId);
	var focused = overlay.data('focus');
	$(document).off("keypress.closeBook");
	maxZ = overlay.data("prevZ");
	overlay.hide();
	if(focused!=null){
		focused.focus();
	}
}

function toggleScraperSelect(){
	var node = $(this);
	var parent = node.parent();
	if(!node.hasClass("selected")){
		$(".selected",parent).removeClass("selected");
	}
	node.toggleClass("selected");
}
function addScraperAuthor(author, parent){
	if(parent==null){
		parent = $("#scraperAuthorTo");
	}
	var row = $("<div></div>").addClass("selectable").data("data",author);
	var toAuthorId = $("<div></div>").addClass("field id");;
	toAuthorId.html(author.id);
	toAuthorId.appendTo(row);
	var toAuthorName = $("<div></div>").addClass("field").attr("name","name");
	toAuthorName.html(author.name);
	toAuthorName.appendTo(row);
	var toAuthorUrl = $("<div></div>").addClass("field url").attr("name","url");
	toAuthorUrl.html(author.url);
	toAuthorUrl.appendTo(row);
	row.appendTo(parent);
	row.on("click",toggleScraperSelect);
}
function addScraperSeries(series, parent){
	if(parent==null){
		parent = $("#scraperSeriesTo");
	}
	var row = $("<div></div>").addClass("selectable").data("data",series);
	var toSeriesId = $("<div></div>").addClass("field id");
	toSeriesId.html(series.id);
	toSeriesId.appendTo(row);
	var toSeriesName = $("<div></div>").addClass("field").attr("name","name");
	toSeriesName.html(series.name);
	toSeriesName.appendTo(row);
	var toSeriesSeqNo = $("<div></div>").addClass("field seqno").attr("name","seqno");
	toSeriesSeqNo.html(series.seqno);
	toSeriesSeqNo.appendTo(row);
	var toSeriesUrl = $("<div></div>").addClass("field url").attr("name","url");
	toSeriesUrl.html(series.url);
	toSeriesUrl.appendTo(row);
	row.appendTo(parent);
	row.on("click",toggleScraperSelect);
}
function addScraperTag(tag, parent){
	if(parent==null){
		parent = $("#scraperTagsTo");
	}
	var row = $("<div></div>").addClass("selectable").data("data",tag);
	var toTagId = $("<div></div>").addClass("field id");
	toTagId.html(tag.id);
	toTagId.appendTo(row);
	var toTag = $("<div></div>").addClass("field").attr("name","name");
	toTag.html(tag.name);
	toTag.appendTo(row);
	row.appendTo(parent);
	row.on("click",toggleScraperSelect);
}

function doPopupScraperBook(metadata){
	var focused = $(':focus');
	var prevZ = maxZ;
	maxZ = maxZ+10;
	var thisZ = maxZ;
	var book = readBackBook();
	var overlay = $('#BookScraperOverlay');
    overlay.data('focus',focused).data('prevZ',prevZ);
    var box =$('.popup',overlay);
	var content=$(".popupContent",box);
	content.data("data",book);
	// local book part
	$("#scraperTitleTo").html(book.title).data("data",book.title);

	var authorsTo = $("#scraperAuthorTo");
	authorsTo.empty();
	$.each(book.authors, function(key, author) {
		addScraperAuthor(author,authorsTo);
	});

	var seriesTo = $("#scraperSeriesTo");
	seriesTo.empty();
	$.each(book.series, function(key, series) {
		addScraperSeries(series,seriesTo);
	});

	var tagsTo = $("#scraperTagsTo");
	tagsTo.empty();
	$.each(book.tags, function(key, tag) {
		addScraperTag(tag, tagsTo)
	});

	$("#scraperBlurbTo").html(book.blurb.replaceAll("\n","<br>")).data("data",book.blurb);;

	// metadata part
	$("#scraperTitleFrom").html(metadata.title).data("data",metadata.title);;

	var authorsFrom = $("#scraperAuthorFrom");
	authorsFrom.empty();
	$.each(metadata.authors, function(key, author) {
		var row = $("<div></div>").addClass("selectable").data("data",author);
		var frmAuthorName = $("<div></div>").addClass("field").attr("name","name");
		frmAuthorName.html(author.name);
		frmAuthorName.appendTo(row);
		var frmAuthorUrl = $("<div></div>").addClass("field url").attr("name","url");
		frmAuthorUrl.html(author.url);
		frmAuthorUrl.appendTo(row);
		row.appendTo(authorsFrom);
		row.on("click",toggleScraperSelect);
	});

	var seriesFrom = $("#scraperSeriesFrom");
	seriesFrom.empty();
	$.each(metadata.series, function(key, series) {
		var row = $("<div></div>").addClass("selectable").data("data",series);
		var frmSeriesName = $("<div></div>").addClass("field").attr("name","name");
		frmSeriesName.html(series.name);
		frmSeriesName.appendTo(row);
		var frmSeriesSeqNo = $("<div></div>").addClass("field seqno").attr("name","seqno");
		frmSeriesSeqNo.html(series.seqno);
		frmSeriesSeqNo.appendTo(row);
		var frmSeriesUrl = $("<div></div>").addClass("field url").attr("name","url");
		frmSeriesUrl.html(series.url);
		frmSeriesUrl.appendTo(row);
		row.appendTo(seriesFrom);
		row.on("click",toggleScraperSelect);
	});

	var tagsFrom = $("#scraperTagsFrom");
	tagsFrom.empty();
	$.each(metadata.tags, function(key, tag) {
		var objTag = {name:tag};
		var frmTag = $("<div></div>").addClass("field selectable").attr("name","name").data("data",objTag);
		frmTag.html(tag);
		frmTag.appendTo(tagsFrom);
		frmTag.on("click",toggleScraperSelect);
	});

	$("#scraperBlurbFrom").html(metadata.blurb.replaceAll("\n","<br>")).data("data",metadata.blurb);


	$(document).on("keypress.closeBook", function(e) {
		if (e.key === "Escape" || e.key==="Enter") {
			overlay.click();
		}
		e.preventDefault();
		e.stopPropagation();
	});
    overlay.show();

	focused.blur();
	overlay.show().on("click",function(e){
		if (e.target != e.currentTarget) return;
		closePopupBook(e,"#BookScraperOverlay")
	});
	resizePopup(box,content);
	//preventScroll(overlay);
	//spreventScroll(box);
}
function editBook(event){
	var overlay = $('#BookOverlay');
	var bookId=overlay.data("bookId");
	overlay.click();
	changePage("BookEdit");
	loadABook(bookId,editABook)
}


function scraperMoveStatic(fromId, toId){
	var from=$(fromId);
	var value = from.html();
	$(toId).html(value).data("data",from.data("data"));
}

function scraperMoveList(fromId, toId){
	var frCell = $(fromId+" .selected");
	var toCell = $(toId+" .selected");
	if((frCell.length==0) || (toCell.length==0)) return
	var fromData = frCell.data("data");
	var toData = toCell.data("data");
	for (var [key, value] of Object.entries(fromData)) {
		toData[key]=value;
		var field = $("div[name='"+key+"']",toCell).html(value);
	}
	toCell.data("data",toData);
}
function scraperAddList(fromId,func){
	var frCell = $(fromId+" .selected");
	if(frCell.length==0) return
	var fromData = frCell.data("data");
	fromData["id"]="";
	func(fromData);
}