import React from 'react';
import MenuItem from './MenuItem';

function MenuBar() {
	return (
		<div className="nav flex-column sidebar bg-secondary text-white text-center">
            <h4 className="mt-3"><strong>DFD</strong></h4>
            <MenuItem menuName="Search" link="/search" />
            <MenuItem menuName="New Topic" link="/new-topic"/>
            <MenuItem menuName="Viewed" link="/viewed" />
            <MenuItem menuName="Settings" link="/settings" />
		</div>
	);
}

export default MenuBar;
