CREATE OR REPLACE PROCEDURE universities_setup() 
LANGUAGE plpgsql
AS $$
DECLARE
 names TEXT[] := ARRAY['Belarussian University of Informatics and Radioelectronics', 'Belarussian state universite', 'Higher school of economics','Moscow Insititute of Physics and Technology', 'Saint-Petersburg National Research University of Information Technologies, Mechanics and Optics'];
 shorthands TEXT[] := ARRAY['BSUIR', 'BSU', 'HSE', 'MIPT', 'ITMO'];
BEGIN 
      FOR i in 1..5 LOOP
            INSERT INTO universities (full_name, shorthand) 
            SELECT names[i], shorthands[i];
     END LOOP; 
END
$$;