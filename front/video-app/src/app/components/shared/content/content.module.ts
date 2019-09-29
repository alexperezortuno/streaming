import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ContentComponent} from './content.component';
import {MatCardModule} from "@angular/material/card";

@NgModule({
    declarations: [ContentComponent],
    imports: [
        CommonModule,
        MatCardModule
    ],
    exports: [
        ContentComponent
    ]
})
export class ContentModule {
}
