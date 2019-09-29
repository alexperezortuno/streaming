import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {DashComponent} from './dash.component';
import {DashRoutingModule} from "./dash-routing.module";
import {NavbarModule} from "../../shared/navbar/navbar.module";
import {ContentModule} from "../../shared/content/content.module";

@NgModule({
    declarations: [DashComponent],
    imports: [
        DashRoutingModule,
        CommonModule,
        ContentModule,
        NavbarModule,
    ]
})
export class DashModule {
}
