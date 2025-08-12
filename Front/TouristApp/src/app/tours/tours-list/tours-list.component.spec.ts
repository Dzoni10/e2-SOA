import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ToursListComponent } from './tours-list.component';

describe('ToursListComponent', () => {
  let component: ToursListComponent;
  let fixture: ComponentFixture<ToursListComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ToursListComponent]
    });
    fixture = TestBed.createComponent(ToursListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
