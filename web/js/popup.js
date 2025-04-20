var maxZ = 2000;


function resizePopup(box,content){
	var origheight = box.outerHeight();
	var maxheight = $(window).height()-80;
	height = Math.min(origheight,maxheight);
	console.debug({
		origheight:origheight,
		maxheight:maxheight,
		height:height
	})
	box.css({
		height:height});
	
	if(content!=null && origheight>maxheight){
		content = $(content).css("overflow-y","auto").height(height);
		//parent = content.;
        /*
		scroll = parent.niceScroll(content,{cursorwidth:'8px',
			cursorborder:'none',
			cursorborderradius:'3px',
			cursorcolor:'#db770e',
			cursoropacitymax:'0.75'
		});
        */
	}
	

}

