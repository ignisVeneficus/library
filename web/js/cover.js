function createBlankCoverContent(parent,textColor, title, author){
	parent.empty();
	var author = $("<div>"+author+"</div>").addClass("coverAuthor").css("color",textColor);
	author.appendTo(parent);
	var title = $("<div>"+title+"</div>").addClass("coverTitle").css("color",textColor);
	title.appendTo(parent);	

}

function createBlankCover(backgroundColor,textColor, title, author){
	var cover = $("<div></div>").addClass("cover").css({backgroundColor:backgroundColor});
    createBlankCoverContent(cover,textColor, title, author);
	return cover;
}
