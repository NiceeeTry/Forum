var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}




const commentForm = document.querySelector('#comment-form');
const commentList = document.querySelector('#comment-list');

commentForm.addEventListener('submit', (e) => {
  e.preventDefault();
  
  const name = commentForm.name.value;
  const email = commentForm.email.value;
  const comment = commentForm.comment.value;
  
  const timestamp = new Date().toLocaleString();
  
  const newComment = document.createElement('div');
  newComment.classList.add('comment');
  
  newComment.innerHTML = `
    <h3>${name}</h3>
    <p>${comment}</p>
    <p class="timestamp">${timestamp}</p>
  `;
  
  commentList.appendChild(newComment);
  
  commentForm.reset();
});


// likes


