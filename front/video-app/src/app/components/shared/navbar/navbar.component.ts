import {Component, OnInit} from '@angular/core';

@Component({
    selector: 'navbar',
    templateUrl: './navbar.component.html',
    styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit {
    sidenav: boolean;
    constructor() {
    }

    ngOnInit() {
    }

    toggle() {
        this.sidenav = !this.sidenav;
        console.log(this.sidenav);
        return this.sidenav;
    }
}
