
function displayBookList(data){
	var page =$("#BookList");
	displayPagination(page,data.pagination,displayBookList);
	/* filter */
	var filterNode = $(".filters",page);
	fillFilters(filterNode,data.filter,loadAllBookList);

	// count
	var countNode = $(".nrOfBooks",page);
	countNode.html(data.pagination.qty);

	// content
	var content = $(".content",page);
	content.empty();
	$.each(data.result, function(key, book) {
		var color = book.coverColor;
		if(color==""){
			color = colors[book.id%colors.length];
		}
		var c = hexToRgb(color);
		var authors = [];
		for( var i=0;i<book.authors.length;i++){
			authors.push(book.authors[i].name);
		}
		var authorName = authors.join("; ");
		var op = "0.2";
		if(book.edited!=0){
			op = "0.4";
		}
		var div = $("<div></div>").attr("id","book_cover_"+book.id).data("id",book.id).addClass("coverPanel").css({background: "rgba( "+c[0]+","+c[1]+","+c[2]+", "+op });
		div.appendTo(content);
		/*
		if(book.edited!=0){
			var checkbox = $("<div class=\"checked fa-check fa-solid\"></div>").appendTo(div);
		}
		*/
		if(book.edited!=0){
			div.addClass("edited");
		}
		var cover ="";
		if(book.hasCover){
			cover = $("<div></div>").css({backgroundImage:"url(cover/"+book.file+"."+book.coverType+")"}).addClass("cover");
		}else{
			cover = createBlankCover(color,setContrast(c),book.title,authorName)
		}
		cover.appendTo(div);

		var author = $("<div>"+authorName+"</div>").addClass("author");
		author.appendTo(div);
		var line = $("<div></div>").addClass("line");
		line.appendTo(div)
		var title = $("<div>"+book.title+"</div>").addClass("title");
		title.appendTo(div);
		if(data.display){
			$("<div></div>").addClass("line").appendTo(div)
			if(data.display.type=="series"){
				$.each(book.series, function(key, series) {
					if(series.id == data.display.data){
						var text = series.name;
						if((series.seqno!=null)&&(series.seqno!=0)){
							text += " #"+series.seqno;
						}
						$("<div>"+text+"</div>").addClass("display").appendTo(div);
					}
				});
			}
		}
		var fv = "";
		if(book.fileType == "mobi"){
			fv = "Mobipocket"
		}
		if(book.fileType == "epub"){
			fv = "Epub"
		}

		var fileType = $("<div class=\"fileformat\">"+fv+"</div>").appendTo(div);
        div.on("click",function(event){
            loadABook(book.id,displayABook)
        });
	});
}

function displayABook(data){
    doPopupBook(data);
}


function addAuthorRow(before,author,addOrig){
	if(before == null){
		before = $("#newAuthorRow");
	}
	var row = $("<div class=\"row\"></div>").data("hasOrig",author!=null&&addOrig);
	row.insertBefore(before);
	var id = $("<div class=\"id\"></div>");
	id.appendTo(row);
	var nameN=$("<input name=\"name\"></input>");
	nameN.appendTo(row);
	var urlN=$("<input name=\"url\"></input>");
	urlN.appendTo(row);
	var btnN = $("<div class=\"buttons\"></div>");
	btnN.appendTo(row);
	var deleteBtn = $("<button class=\"button formButton fa-solid fa-trash-can\" title=\"Remove author\"></button>").addTooltip();
	deleteBtn.appendTo(btnN);
	deleteBtn.on("click",function(event){
		$(this).parent().parent().remove();
	});
	if(author!=null){
		id.html(author.id);
		nameN.val(author.name);
		urlN.val(author.url);
		if(addOrig){
			nameN.data("orig",author.name).on( "focusout",checkRow)
			urlN.data("orig",author.url).on( "focusout",checkRow)
		}
	}

}
function fillEditAuthors(authors){
	var node = $("#BookEdit_Authors");
	node.data("orig",authors);
	node.empty();
	var row = $("<div id=\"newAuthorRow\" class=\"row\"></div>");
	row.appendTo(node);
	// hack the ccs, empty eleme for the first, so button go to end
	var btnN = $("<div class=\"buttons\"></div>");
	btnN.appendTo(row);
	btnN = $("<div class=\"buttons\"></div>");
	btnN.appendTo(row);
	var addBtn = $("<button class=\"button formButton fa-solid fa-plus\" title=\"Add new author\"></button>").addTooltip();
	addBtn.appendTo(btnN);
	addBtn.on("click",function(event){
		addAuthorRow(row);
	});

	$.each(authors, function(key, author) {
		addAuthorRow(row,author,true);
	});

}
function addSeriesRow(before,series,addOrig){
	if(before==null){
		before = $("#newSeriesRow");
	}
	var row = $("<div class=\"row\"></div>");
	row.insertBefore(before);
	var id = $("<div class=\"id\"></div>");
	id.appendTo(row);
	var nameN=$("<input name=\"name\"></input>");
	nameN.appendTo(row);
	var seqN=$("<input name=\"seqno\" class=\"seqno\"></input>");
	seqN.appendTo(row);
	var urlN=$("<input name=\"url\"></input>");
	urlN.appendTo(row);
	var btnN = $("<div class=\"buttons\"></div>");
	btnN.appendTo(row);
	var deleteBtn = $("<button class=\"button formButton fa-solid fa-trash-can\" title=\"Remove series\"></button>").addTooltip();
	deleteBtn.appendTo(btnN);
	deleteBtn.on("click",function(event){
		$(this).parent().parent().remove();
	});
	if(series!=null){
		id.html(series.id);
		nameN.val(series.name);
		seqN.val(series.seqno);
		urlN.val(series.url);
		if(addOrig){
			nameN.data("orig",series.name).on( "focusout",checkRow);
			seqN.data("orig",series.seqno).on( "focusout",checkRow);
			urlN.data("orig",series.url).on( "focusout",checkRow);
		}
	}

}
function fillEditSeries(series){
	node = $("#BookEdit_Series");
	node.data("orig",series);
	node.empty();

	var row = $("<div id=\"newSeriesRow\" class=\"row\"></div>");
	row.appendTo(node);
	btnN = $("<div class=\"buttons\"></div>");
	btnN.appendTo(row);
	btnN = $("<div class=\"buttons\"></div>");
	btnN.appendTo(row);
	addBtn = $("<button class=\"button formButton fa-solid fa-plus\" title=\"Add new series\"></button>").addTooltip();
	addBtn.appendTo(btnN);
	addBtn.on("click",function(event){
		addSeriesRow(row);
	});

	$.each(series, function(key, se) {
		addSeriesRow(row,se,true)
	});

}

function addTagRow(before,tag,addOrig){
	if(before==null){
		before = $("#newTagRow");
	}
	var row = $("<div class=\"row\"></div>");
	row.insertBefore(before);
	var id = $("<div class=\"id\"></div>");
	id.appendTo(row);
	var nameN=$("<input name=\"name\"></input>");
	nameN.appendTo(row);
	var btnN = $("<div class=\"buttons\"></div>")
	btnN.appendTo(row);
	var deleteBtn = $("<button class=\"button formButton fa-solid fa-trash-can\" title=\"Remove tag\"></button>").addTooltip();
	deleteBtn.appendTo(btnN);
	deleteBtn.on("click",function(event){
		$(this).parent().parent().remove();
	});
	if(tag!=null){
		id.html(tag.id);
		nameN.val(tag.name);
		if(addOrig){
			nameN.data("orig",tag.name).on( "focusout",checkRow);
		}
	}

}
function fillEditTags(tags){
	node = $("#BookEdit_Tags");
	node.data("orig",tags);
	node.empty();
	
	var row = $("<div id=\"newTagRow\" class=\"row\"></div>");
	row.appendTo(node);
	btnN = $("<div class=\"buttons\"></div>");
	btnN.appendTo(row);
	btnN = $("<div class=\"buttons\"></div>");
	btnN.appendTo(row);
	addBtn = $("<button class=\"button formButton fa-solid fa-plus\" title=\"Add new tag\"></button>").addTooltip();
	addBtn.appendTo(btnN);
	addBtn.on("click",function(event){
		addTagRow(row);
	});

	$.each(tags, function(key, tag) {
		addTagRow(row,tag,true)
	});

}
function checkRow(){
	checkRowWith($(this).parent());
}
function checkRowWith(row){
	if(row==null) row = $(this).parent();
	var needCheck = row.data("hasOrig");
	if(!needCheck) return;

	var input = $("input",row)
	var hasChanged = false;
	input.each(function(){
		if($(this).val()!=$(this).data("orig")){
			hasChanged=true
			$(this).addClass("changed");
		}
	});
	if(!hasChanged){
		$(".warning",row).remove();
		$(".changed",row).removeClass("changed");
		return;
	} 
	var warning = $(".warning",row);
	if(warning.length!=0){
		return;
	}
	warning = $("<button class=\"warning button formButton fa-solid fa-rotate-left\" title=\"Discard changes\"></button>").on("click", function(event){
		var row = $(this).parent().parent();
		var input = $("input",row);
		input.each(function(){
			$(this).val($(this).data("orig"));
		});
		checkRowWith(row);
	}).addTooltip();
	$(".buttons",row).prepend(warning);
}

function readBackBookComponent(node){
	var comp = [];
	var node = $(node);
	var rows = $(".row",node);
	rows.each(function(){
		var row = $(this);
		inputs = $("input",row)
		if(inputs.length!=0){
			var data = {};
			data.id = $(".id",row).html();
			inputs.each(function(){
				input = $(this);
				data[input.attr("name")]=input.val()
			});
			comp.push(data);
		}
	});
	return comp;

}

function readBackBook(){
	var root = $("#BookEdit .content");
	orig = root.data("orig");
	book = {};
	book.title= $("#BookEdit_title").val();
	book.url= $("#BookEdit_url").val();
	book.blurb= $("#BookEdit_blurb").val();
	book.id = orig.id;
	book.hasCover = orig.hasCover;
    book.coverColor = orig.coverColor;
	book.coverType = orig.coverType;
    book.file = orig.file;
	book.authors = readBackBookComponent("#BookEdit_Authors");
	book.series = readBackBookComponent("#BookEdit_Series");
	book.tags = readBackBookComponent("#BookEdit_Tags");

	return book;

}

function editABook(book){
	var root = $("#BookEdit .content");
	root.data("orig",book);
    var coverNode =  $(".cover",root);
	$(".warning",root).remove();
	$(".changed",root).removeClass("changed");
	var authors = [];
	for( var i=0;i<book.authors.length;i++){
		authors.push(book.authors[i].name);
	}
	var authorName = authors.join("; ");
	if(book.hasCover){
       coverNode.css({backgroundImage:"url(cover/"+book.file+"."+book.coverType+")",backgroundColor:""})
       coverNode.empty();
    }
    else{
        var color = colors[book.id%colors.length];
        coverNode.css({backgroundImage:"",backgroundColor:color})
        var c = hexToRgb(color);
        createBlankCoverContent(coverNode,setContrast(c),book.title, authorName);
    }

	var titleN = $("#BookEdit_title").val(book.title).data("orig",book.title).on( "focusout",checkRow);
	var urlN = $("#BookEdit_url").val(book.url).data("orig",book.url).on( "focusout",checkRow);
	var blurbN= $("#BookEdit_blurb").val(book.blurb).data("orig",book.blurb).on("focusout",checkRow);

	titleN.parent().data("hasOrig",true);
	urlN.parent().data("hasOrig",true);
	blurbN.parent().data("hasOrig",true);

	fillEditAuthors(book.authors);
	fillEditSeries(book.series);
	fillEditTags(book.tags);
}

function scrapeBook(){
	var url = $("#BookEdit_url").val();
	scrapeABook(url);
}

function writeBackScraper(){
	var content = $("#BookScraperOverlay .popupContent");
	var book = content.data("data");
	book.title=$("#scraperTitleTo").data("data");
	var authors = [];
	$("#scraperAuthorTo .selectable").each(function(idx){
		authors.push($(this).data("data"));
	});
	book.authors = authors;
	var series = [];
	$("#scraperSeriesTo .selectable").each(function(idx){
		series.push($(this).data("data"));
	});
	book.series = series;
	var tags = [];
	$("#scraperTagsTo .selectable").each(function(idx){
		tags.push($(this).data("data"));
	});
	book.tags = tags;
	book.blurb=$("#scraperBlurbTo").data("data");

	$("#BookScraperOverlay").click();

	replaceEditedBookData(book)
}

function replaceEditedBookData(book){
	var root = $("#BookEdit .content");
 	$(".warning",root).remove();
	$(".changed",root).removeClass("changed");
 
	var titleN = $("#BookEdit_title").val(book.title);
	var urlN = $("#BookEdit_url").val(book.url);
	var blurbN= $("#BookEdit_blurb").val(book.blurb);

	checkRowWith(titleN.parent());
	checkRowWith(urlN.parent());
	checkRowWith(blurbN.parent());

	replaceEditData(book.authors,"#BookEdit_Authors",addAuthorRow);
	replaceEditData(book.series,"#BookEdit_Series",addSeriesRow);
	replaceEditData(book.tags,"#BookEdit_Tags",addTagRow);


}
function replaceEditData(data,parent,addnew){
	var node = $(parent);
	var rows = $(".row",node).has( "input" );
	qty = rows.length;
	rows.each(function(index){
		row = $(this);
		datarow = data[index];
		$("input",row).each(function(ri,val){
			node = $(val);
			nname = node.attr("name");
			node.val(datarow[nname]);
		})
		checkRowWith(row);
	});
	if(qty<data.length){
		for(var i=qty;i<data.length;i++){
			addnew(null,data[i],false);
		}
	}

}
