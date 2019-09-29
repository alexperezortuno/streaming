import {NgModule} from '@angular/core';
import {NavbarComponent} from "./navbar.component";
import {MatToolbarModule} from "@angular/material/toolbar";
import {RouterModule} from "@angular/router";
import {MatIconModule} from "@angular/material/icon";
import {MatButtonModule} from "@angular/material/button";

@NgModule({
    declarations: [
        NavbarComponent
    ],
    imports: [
        MatToolbarModule,
        RouterModule,
        MatIconModule,
        MatButtonModule
    ],
    exports: [
        NavbarComponent
    ]
})
export class NavbarModule {
}
