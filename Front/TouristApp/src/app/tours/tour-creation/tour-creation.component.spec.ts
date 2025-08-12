import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TourCreationComponent } from './tour-creation.component';

describe('TourCreationComponent', () => {
  let component: TourCreationComponent;
  let fixture: ComponentFixture<TourCreationComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [TourCreationComponent]
    });
    fixture = TestBed.createComponent(TourCreationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
