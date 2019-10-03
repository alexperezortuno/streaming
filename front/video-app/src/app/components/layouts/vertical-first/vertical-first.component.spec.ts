import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VerticalFirstComponent } from './vertical-first.component';

describe('VerticalFirstComponent', () => {
  let component: VerticalFirstComponent;
  let fixture: ComponentFixture<VerticalFirstComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VerticalFirstComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VerticalFirstComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
